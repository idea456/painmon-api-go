package types

import model "github.com/idea456/painmon-api-go/graph/model"

func ToPointers[T interface{}](slice []T) []*T {
	arr := make([]*T, 0)

	for item := range slice {
		arr = append(arr, &slice[item])
	}

	return arr
}

func MapDomain(domain Domain) *model.Domain {
	return &model.Domain{
		ID:    domain.Name,
		Name:  domain.Name,
		Ar:    &domain.Level,
		Level: &domain.Level,
	}
}

func MapDomainCategory(domainCategory DomainCategory) model.DomainCategory {
	domains := make([]*model.Domain, 0)
	for _, domain := range domainCategory.Domains {
		domains = append(domains, MapDomain(domain))
	}

	return model.DomainCategory{
		Name:      domainCategory.Name,
		Domains:   domains,
		Artifacts: ToPointers[string](domainCategory.Artifacts),
	}
}
