package dnsbl

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

func (r *mutationResolver) Enqueue(ctx context.Context, ipAddresses []string) (string, error) {
	err := Enqueue(ctx, r.db, ipAddresses)
	return "", err
}

func (r *queryResolver) GetIPDetails(ctx context.Context, ipAddress string) (*IPDetail, error) {
	return GetIPDetail(ctx, r.db, ipAddress)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
