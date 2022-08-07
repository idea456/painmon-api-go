package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/idea456/painmon-api-go/graph/generated"
	"github.com/idea456/painmon-api-go/graph/model"
	"github.com/idea456/painmon-api-go/internal/database"
	"github.com/idea456/painmon-api-go/internal/types"
)

// GetDomainCategories is the resolver for the getDomainCategories field.
func (r *queryResolver) GetDomainCategories(ctx context.Context) ([]*model.DomainCategory, error) {
	domainCategories := make([]*model.DomainCategory, 0)

	domainCategories = append(domainCategories, &model.DomainCategory{
		Name: database.Get[types.Talent]("kamisatoayaka").Name,
		Domains: []*model.Domain{&model.Domain{
			ID:   "test",
			Name: "test",
		}},
	})

	return domainCategories, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Domain(ctx context.Context) (*model.Domain, error) {
	panic(fmt.Errorf("not implemented"))
}
