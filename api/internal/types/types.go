package types

import "time"

// TODO: Include this type after "materials" data is fetched
// type Item struct {
// 	Name        string   `json:"name"`
// 	Description string   `json:"description"`
// 	Rarity      int      `json:"rarity"`
// 	Day         []string `json:"daysofweek"`
// 	Type        string   `json:"materialtype"`
// 	Domain      string   `json:"dropdomain"`
// }

type Talent struct {
	Name  string     `json:"name"`
	Costs TalentCost `json:"costs"`
}

type TalentCost struct {
	Cost []map[string]string `json:"lvl10"`
}

// type TalentMaterial struct {
// 	Name          string   `json:"name"`
// 	TwoStarName   string   `json:"2starname"`
// 	ThreeStarName string   `json:"3starname"`
// 	FourStarName  string   `json:"4starname"`
// 	Day           []string `json:"day"`
// 	Domain        string   `json:"domainofmastery"`
// }

type Artifact struct {
	Name   string   `json:"name"`
	Rarity []string `json:"rarity"`
}

type Material struct {
	Name          string `json:"name"`
	Type          string
	TwoStarName   string   `json:"2starname"`
	ThreeStarName string   `json:"3starname"`
	FourStarName  string   `json:"4starname"`
	Day           []string `json:"day"`
	Domain        string   `json:"domainofmastery"`
}

type Weapon struct {
	Name     string `json:"name"`
	Rarity   string `json:"rarity"`
	Material string `json:"weaponmaterialtype"`
}

// type WeaponMaterial struct {
// 	Name          string   `json:"name"`
// 	TwoStarName   string   `json:"2starname"`
// 	ThreeStarName string   `json:"3starname"`
// 	FourStarName  string   `json:"4starname"`
// 	Day           []string `json:"day"`
// 	Domain        string   `json:"domainofforgery"`
// }

type Domain struct {
	Name           string   `json:"name"`
	Region         string   `json:"region"`
	DomainType     string   `json:"domaintype"`
	DomainCategory string   `json:"domainentrance"`
	Level          int      `json:"recommendedlevel"`
	Rewards        []Reward `json:"rewardpreview"`
	Monsters       []string `json:"monsterlist"`
	Disorder       []string `json:"disorder"`
}

type DomainCategory struct {
	Name      string
	Domains   []Domain
	Artifacts []string
}

type Reward struct {
	Name   string `json:"name"`
	Count  int    `json:"count"`
	Rarity string `json:"rarity"`
}

type Entry interface {
	Artifact | Talent | Weapon | Material | Domain | DomainCategory
}
type Daily struct {
	Date      time.Time
	Day       string
	Materials []Material
	// TalentMaterials []Material
	// WeaponMaterials []Material
}
