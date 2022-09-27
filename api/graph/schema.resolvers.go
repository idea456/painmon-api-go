package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"sync"

	"github.com/idea456/painmon-api-go/graph/generated"
	"github.com/idea456/painmon-api-go/graph/model"
	"github.com/idea456/painmon-api-go/internal/controllers"
	"github.com/idea456/painmon-api-go/internal/database"
	"github.com/idea456/painmon-api-go/internal/types"
	"github.com/idea456/painmon-api-go/pkg/utils"
)

// GetDomainCategories is the resolver for the getDomainCategories field.
func (r *queryResolver) GetDomainCategories(ctx context.Context) ([]*model.DomainCategory, error) {
	domainCategories := make([]*model.DomainCategory, 0)
	categories := make(map[string][]types.Domain)

	var category map[string]types.Domain = database.GetCategory[types.Domain]("domains")

	for domainName := range category {
		domain := category[domainName]

		if _, ok := categories[domain.DomainCategory]; !ok {
			categories[domain.DomainCategory] = make([]types.Domain, 0)
		}
		categories[domain.DomainCategory] = append(categories[domain.DomainCategory], domain)
	}

	var wg sync.WaitGroup
	wg.Add(len(categories))

	for domainCategory := range categories {
		go func(domainCategory string) {
			defer wg.Done()

			domains := make([]*model.Domain, 0)
			artifacts := make([]string, 0)

			for _, domain := range categories[domainCategory] {
				domains = append(domains, types.MapDomain(domain))
				for _, artifact := range domain.Rewards {
					if !utils.In(artifacts, artifact.Name) {
						artifacts = append(artifacts, artifact.Name)
					}
				}
			}

			domainCategories = append(domainCategories, &model.DomainCategory{
				Name:      domainCategory,
				Domains:   domains,
				Artifacts: types.ToPointers[string](artifacts),
			})
		}(domainCategory)
	}

	wg.Wait()

	return domainCategories, nil
}

// GetDaily is the resolver for the getDaily field.
func (r *queryResolver) GetDaily(ctx context.Context) (*model.Daily, error) {
	daily := controllers.NewDailyController()

	daily.GetDailies()

	return &model.Daily{
		Date:       daily.Date,
		Day:        &daily.Day,
		Materials:  daily.Materials,
		Characters: daily.Characters,
		Weapons:    daily.Weapons,
	}, nil
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
