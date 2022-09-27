package controllers

import (
	"strconv"
	"sync"
	"time"

	"github.com/idea456/painmon-api-go/graph/model"
	"github.com/idea456/painmon-api-go/internal/database"
	"github.com/idea456/painmon-api-go/internal/types"
	"github.com/idea456/painmon-api-go/pkg/utils"
)

type DailyController struct {
	Date       *time.Time
	Day        string
	Characters []*model.Character
	Materials  []*model.ItemGroup
	Weapons    []*model.Weapon
	Image      string
	Wg         *sync.WaitGroup
}

func NewDailyController() *DailyController {
	today := time.Now()
	day := today.Weekday().String()

	return &DailyController{
		Date:      &today,
		Day:       day,
		Wg:        &sync.WaitGroup{},
		Materials: make([]*model.ItemGroup, 0),
	}
}

func (dc *DailyController) GetDailies() {
	dc.Wg.Add(2)
	dc.GetCharactersAndWeapons()
	dc.GetMaterials()
	dc.Wg.Wait()
}

func (dc *DailyController) GetCharactersAndWeapons() {
	defer dc.Wg.Done()
	daily := database.GetPath[[]types.Material]("daily", ".Tuesday")
	characters := make([]*model.Character, 0)
	weapons := make([]*model.Weapon, 0)

	getRarity := func(rarityStr string) *int {
		rarity, _ := strconv.Atoi(rarityStr)
		return &rarity
	}

	for _, material := range daily {
		switch material.Type {
		case utils.TALENT_MATERIAL_TYPE:
			for i := 0; i < len(material.Characters); i++ {
				character := material.Characters[i]
				characters = append(characters, &model.Character{
					Name:    character.Name,
					Rarity:  getRarity(character.Rarity),
					Element: &character.Element,
					Sex:     &character.Sex,
					Nation:  &character.Nation,
				})
			}
		case utils.WEAPON_MATERIAL_TYPE:
			for i := 0; i < len(material.Weapons); i++ {
				weapon := material.Weapons[i]
				weapons = append(weapons, &model.Weapon{
					Name:        weapon.Name,
					Description: &weapon.Description,
					Type:        &weapon.Type,
					Rarity:      getRarity(weapon.Rarity),
					Atk:         &weapon.Atk,
					Substat:     &weapon.Substat,
					Effectname:  &weapon.EffectName,
					Effect:      &weapon.Effect,
				})
			}
		}

	}

	dc.Characters = characters
	dc.Weapons = weapons
}

func (dc *DailyController) GetMaterials() {
	defer dc.Wg.Done()

	// redis single-threaded no goroutines T_T
	talentMaterials := database.GetCategory[types.Material](utils.TALENT_MATERIAL)
	weaponMaterials := database.GetCategory[types.Material](utils.WEAPON_MATERIAL)

	var mu sync.Mutex

	materials := make([]*model.ItemGroup, 0)

	// NOTE: Don't use the global dc.Wg since we want to separate GetDailies's 3 running goroutines waiting group counter from our filter ones
	var filterWg sync.WaitGroup
	filterWg.Add(2)
	filterMaterial := func(array map[string]types.Material) {
		defer filterWg.Done()

		for key := range array {
			material := array[key]
			if utils.In(material.Day, dc.Day) {
				items := []*model.Item{
					&model.Item{
						ID:   &material.TwoStarName,
						Name: &material.TwoStarName,
					},
					&model.Item{
						ID:   &material.ThreeStarName,
						Name: &material.ThreeStarName,
					},
					&model.Item{
						ID:   &material.FourStarName,
						Name: &material.FourStarName,
					},
				}

				mu.Lock()
				materials = append(materials, &model.ItemGroup{
					Name:   material.Name,
					Day:    &dc.Day,
					Domain: &material.Domain,
					Items:  items,
					Type:   &material.Type,
				})
				mu.Unlock()
			}
		}
	}

	go filterMaterial(talentMaterials)
	go filterMaterial(weaponMaterials)
	filterWg.Wait()

	dc.Materials = materials
}
