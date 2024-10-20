package server

import (
	"context"
	secureRand "crypto/rand"
	"fmt"
	insecureRand "math/rand/v2"

	"math/rand/v2"

	"cloud.google.com/go/firestore"
)

type randomDocumentGetterImpl[T any] struct {
	firestoreClient *firestore.Client
	random          *insecureRand.Rand
	query           firestore.Query
}

func NewRandomDocumentGetter[T any](
	firestoreClient *firestore.Client,
	query firestore.Query,
) (*randomDocumentGetterImpl[T], error) {
	var seedBytes [32]byte
	if _, err := secureRand.Read(seedBytes[:]); err != nil {
		return nil, fmt.Errorf("failed to seed random generator: %w", err)
	}

	return &randomDocumentGetterImpl[T]{
		firestoreClient: firestoreClient,
		random:          rand.New(insecureRand.NewChaCha8(seedBytes)),
		query:           query,
	}, nil
}

func (r *randomDocumentGetterImpl[T]) GetRandomDocuments(
	ctx context.Context,
	limit int,
) ([]T, []*firestore.DocumentSnapshot, error) {
	docs, err := r.query.Documents(ctx).GetAll()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get documents: %w", err)
	}
	if len(docs) == 0 {
		return nil, nil, fmt.Errorf("no documents found")
	}
	if len(docs) < limit {
		return nil, nil, fmt.Errorf("not enough documents found")
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
