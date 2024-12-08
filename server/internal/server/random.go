package server

import (
	"context"
	secureRand "crypto/rand"
	"fmt"
	insecureRand "math/rand/v2"
	"sync/atomic"
	"time"

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
		random:          insecureRand.New(insecureRand.NewChaCha8(seedBytes)),
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
		docs = append([]*firestore.DocumentSnapshot(nil), cached.response...)
	} else {
		docs, err = r.query.Documents(ctx).GetAll()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get documents: %w", err)
		}
		if r.cacheTime > 0 && docs != nil {
			r.cache.Store(&documentCache{
				response:  docs,
				timestamp: time.Now(),
			})
		}
	}
	// TODO(rbtz): generalize to weight.
	var activeDocs []*firestore.DocumentSnapshot
	for _, doc := range docs {
		active_raw := doc.Data()["active"]
		if active_raw == nil {
			continue
		}
		if active := active_raw.(bool); active {
			activeDocs = append(activeDocs, doc)
		}
	}
	if activeDocs == nil || len(activeDocs) == 0 {
		return nil, nil, fmt.Errorf("no documents found")
	}
	if len(activeDocs) < limit {
		return nil, nil, fmt.Errorf("not enough documents found")
	}

	r.random.Shuffle(len(activeDocs), func(i, j int) {
		activeDocs[i], activeDocs[j] = activeDocs[j], activeDocs[i]
	})
	if len(docs) > limit {
		activeDocs = activeDocs[:limit]
	}

	objs := make([]T, len(activeDocs))
	for i, doc := range activeDocs {
		if err := doc.DataTo(&objs[i]); err != nil {
			return nil, nil, fmt.Errorf("failed to decode document: %w", err)
		}
	}

	return objs, activeDocs, nil
}
