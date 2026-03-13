package menus

import (
	"context"
	"strings"

	"github.com/immick00/sako-api/db"
)

type Menus struct {
	Mcdonalds []db.Mcdonald  `json:"mcdonalds,omitempty"`
	Subway    []db.Subway    `json:"subway,omitempty"`
	ChickFilA []db.Chickfila `json:"chick-fil-a,omitempty"`
}

type MenuService struct {
	queries *db.Queries
	ctx     context.Context
}

func New(queries *db.Queries) *MenuService {
	return &MenuService{queries: queries, ctx: context.Background()}
}

func (s *MenuService) GetMenus(ctx context.Context, restaurants []string) (Menus, error) {
	menus := Menus{}

	for _, restaurant := range restaurants {
		switch strings.ToLower(restaurant) {
		case "mcdonalds":
			result, err := s.queries.GetMcdonaldsProducts(ctx)
			if err == nil {
				menus.Mcdonalds = result
			}
		case "subway":
			result, err := s.queries.GetSubwayProducts(ctx)
			if err == nil {
				menus.Subway = result
			}
		case "chick-fil-a":
			result, err := s.queries.GetChickFilAProducts(ctx)
			if err == nil {
				menus.ChickFilA = result
			}
		}
	}

	return menus, nil
}
