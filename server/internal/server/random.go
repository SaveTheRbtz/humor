package server

import (
	"context"
	secureRand "crypto/rand"
	"fmt"
	insecureRand "math/rand/v2"
	"sync/atomic"
	"time"

	"math/rand/v2"

	"cloud.google.com/go/firestore"
)

type documentCache struct {
	response  []*firestore.DocumentSnapshot
	timestamp time.Time
}

type randomDocumentGetterImpl[T any] struct {
	firestoreClient *firestore.Client
	random          *insecureRand.Rand
	query           firestore.Query
	cacheTime       time.Duration
	cache           atomic.Pointer[documentCache]
}

func NewRandomDocumentGetter[T any](
	firestoreClient *firestore.Client,
	query firestore.Query,
	cacheTime time.Duration,
) (*randomDocumentGetterImpl[T], error) {
	var seedBytes [32]byte
	if _, err := secureRand.Read(seedBytes[:]); err != nil {
		return nil, fmt.Errorf("failed to seed random generator: %w", err)
	}

	return &randomDocumentGetterImpl[T]{
		firestoreClient: firestoreClient,
		random:          rand.New(insecureRand.NewChaCha8(seedBytes)),
		query:           query,
		cacheTime:       cacheTime,
	}, nil
}

func (r *randomDocumentGetterImpl[T]) GetRandomDocuments(
	ctx context.Context,
	limit int,
) ([]T, []*firestore.DocumentSnapshot, error) {
	var docs []*firestore.DocumentSnapshot
	var err error

	if cached := r.cache.Load(); cached != nil && cached.timestamp.Add(r.cacheTime).After(time.Now()) {
		docs = cached.response[:]
	} else {
		docs, err = r.query.Documents(ctx).GetAll()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get documents: %w", err)
		}
	}

	if len(docs) == 0 {
		return nil, nil, fmt.Errorf("no documents found")
	}
	if len(docs) < limit {
		return nil, nil, fmt.Errorf("not enough documents found")
	}

	if r.cacheTime > 0 {
		r.cache.Store(&documentCache{
			response:  docs,
			timestamp: time.Now(),
		})
	}

	r.random.Shuffle(len(docs), func(i, j int) {
		docs[i], docs[j] = docs[j], docs[i]
	})
	if len(docs) > limit {
		docs = docs[:limit]
	}

	objs := make([]T, len(docs))
	for i, doc := range docs {
		if err := doc.DataTo(&objs[i]); err != nil {
			return nil, nil, fmt.Errorf("failed to decode document: %w", err)
		}
	}

	return objs, docs, nil
}
