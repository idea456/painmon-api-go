package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/idea456/painmon-api-go/graph/generated"
	model "github.com/idea456/painmon-api-go/graph/model"
	"github.com/idea456/painmon-api-go/internal/database"
	"github.com/idea456/painmon-api-go/internal/types"
	utils "github.com/idea456/painmon-api-go/pkg/utils"
)

// GetDomainCategories is the resolver for the getDomainCategories field.
func (r *queryResolver) GetDomainCategories(ctx context.Context) ([]*model.DomainCategory, error) {
	domainCategories := make([]*model.DomainCategory, 0)
	categories := make(map[string][]types.Domain)

	var category map[string]types.Domain = database.GetCategory[types.Domain]("domains")

	// var category map[string]types.Domain
	// json.Unmarshal(database.GetCategory("domains"), &category)

	for domainName := range category {
		domain := category[domainName]

		if _, ok := categories[domain.DomainCategory]; !ok {
			categories[domain.DomainCategory] = make([]types.Domain, 0)
		}
		categories[domain.DomainCategory] = append(categories[domain.DomainCategory], domain)
	}

	for domainCategory := range categories {
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
	}

	return domainCategories, nil
}

// GetDaily is the resolver for the getDaily field.

// type Daily {
//     date: Time
//     day: String
//     materials: [ItemGroup]
// }
func (r *queryResolver) GetDaily(ctx context.Context) (*model.Daily, error) {
	today := time.Now()
	day := today.Weekday().String()

	talentMaterials := database.GetCategory[types.TalentMaterial](utils.TALENT_MATERIAL)
	weaponMaterials := database.GetCategory[types.WeaponMaterial](utils.WEAPON_MATERIAL)
	materials := make([]*model.ItemGroup, 0)

	for key := range talentMaterials {
		talentMaterial := talentMaterials[key]
		if utils.In(talentMaterial.Day, day) {
			materials = append(materials, &model.ItemGroup{
				Name:   talentMaterial.Name,
				Day:    &day,
				Domain: &talentMaterial.Domain,
				Items:  make([]*model.Item, 0),
				Type:   &utils.TALENT_MATERIAL_TYPE,
			})
		}
	}

	for key := range weaponMaterials {
		weaponMaterial := weaponMaterials[key]
		if utils.In(weaponMaterial.Day, day) {
			materials = append(materials, &model.ItemGroup{
				Name:   weaponMaterial.Name,
				Day:    &day,
				Domain: &weaponMaterial.Domain,
				Items:  make([]*model.Item, 0),
				Type:   &utils.WEAPON_MATERIAL_TYPE,
			})
		}
	}

	return &model.Daily{
		Date:      &today,
		Day:       &day,
		Materials: materials,
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
