package menus

import (
	"context"
	"strings"

	"github.com/immick00/sako-api/db"
)

type Menus struct {
	Mcdonalds  []db.Mcdonald   `json:"mcdonalds,omitempty"`
	Subway     []db.Subway     `json:"subway,omitempty"`
	ChickFilA  []db.Chickfila  `json:"chickfila,omitempty"`
	BurgerKing []db.Burgerking `json:"burgerking,omitempty"`
	TacoBell   []db.Tacobell   `json:"tacobell,omitempty"`
	Popeyes    []db.Popeye     `json:"popeyes,omitempty"`
	Wendys     []db.Wendy      `json:"wendys,omitempty"`
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

	type entry struct {
		keyword string
		fetch   func() error
	}

	entries := []entry{
		{"mcdonalds", func() error {
			result, err := s.queries.GetMcdonaldsProducts(ctx)
			if err == nil {
				menus.Mcdonalds = result
			}
			return err
		}},
		{"subway", func() error {
			result, err := s.queries.GetSubwayProducts(ctx)
			if err == nil {
				menus.Subway = result
			}
			return err
		}},
		{"chickfila", func() error {
			result, err := s.queries.GetChickFilAProducts(ctx)
			if err == nil {
				menus.ChickFilA = result
			}
			return err
		}},
		{"burgerking", func() error {
			result, err := s.queries.GetBurgerKingProducts(ctx)
			if err == nil {
				menus.BurgerKing = result
			}
			return err
		}},
		{"tacobell", func() error {
			result, err := s.queries.GetTacoBellProducts(ctx)
			if err == nil {
				menus.TacoBell = result
			}
			return err
		}},
		{"popeyes", func() error {
			result, err := s.queries.GetPopeyesProducts(ctx)
			if err == nil {
				menus.Popeyes = result
			}
			return err
		}},
		{"wendys", func() error {
			result, err := s.queries.GetWentdysProducts(ctx)
			if err == nil {
				menus.Wendys = result
			}
			return err
		}},
	}

	for _, restaurant := range restaurants {
		normalized := strings.ReplaceAll(strings.ToLower(restaurant), " ", "")
		normalized = strings.ReplaceAll(normalized, "-", "")
		for _, e := range entries {
			if strings.Contains(normalized, e.keyword) || strings.Contains(e.keyword, normalized) {
				e.fetch()
				break
			}
		}
	}

	return menus, nil
}

func (s *MenuService) GetMenusOld(ctx context.Context, restaurants []string) (Menus, error) {
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
		case "burgerking":
			result, err := s.queries.GetBurgerKingProducts(ctx)
			if err == nil {
				menus.BurgerKing = result
			}
		case "tacobell", "taco bell", "taco-bell":
			result, err := s.queries.GetTacoBellProducts(ctx)
			if err == nil {
				menus.TacoBell = result
			}
		case "popeyes", "popeyes louisiana kitchen":
			result, err := s.queries.GetPopeyesProducts(ctx)
			if err == nil {
				menus.Popeyes = result
			}
		}
	}

	return menus, nil
}
