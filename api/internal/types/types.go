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

type Character struct {
	Name        string `json:"name"`
	Book        string
	Day         string
	Description string    `json:"description"`
	Rarity      string    `json:"rarity"`
	Element     string    `json:"element"`
	Sex         string    `json:"gender"`
	Nation      string    `json:"region"`
	Affiliation string    `json:"affiliation"`
	Birthday    string    `json:"birthday"`
	Ascension   Ascension `json:"costs"`
	Materials   Talent
}

type Ascension struct {
	Costs []Cost `json:"ascend6"`
}
type Cost struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}
type Talent struct {
	Name  string     `json:"name"`
	Costs TalentCost `json:"costs"`
}

type TalentCost struct {
	Cost []Cost `json:"lvl10"`
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
	Characters    []Character
	Weapons       []Weapon
}

type Weapon struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"weapontype"`
	Day         string
	Atk         int    `json:"baseatk"`
	Substat     string `json:"substat"`
	EffectName  string `json:"effectname"`
	Effect      string `json:"effect"`
	Rarity      string `json:"rarity"`
	Material    string `json:"weaponmaterialtype"`
}

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
	Artifact | Character | Talent | Weapon | Material | Domain | DomainCategory | []Material
}
type Daily struct {
	Date       time.Time
	Day        string
	Materials  []Material
	Characters []Character
}
