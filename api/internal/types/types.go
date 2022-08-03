package types

type Talent struct {
	Name  string     `json:"name"`
	Costs TalentCost `json:"costs"`
}

type TalentCost struct {
	Cost []map[string]string `json:"lvl10"`
}

type TalentMaterial struct {
	Name   string   `json:"4starname"`
	Day    []string `json:"day"`
	Domain string   `json:"domainofmastery"`
}

type Artifact struct {
	Name   string   `json:"name"`
	Rarity []string `json:"rarity"`
}

type Weapon struct {
	Name     string `json:"name"`
	Rarity   string `json:"rarity"`
	Material string `json:"weaponmaterialtype"`
}
type WeaponMaterial struct {
	Name   string   `json:"name"`
	Day    []string `json:"day"`
	Domain string   `json:"domainofforgery"`
}

type Entry interface {
	Artifact | Talent | TalentMaterial | Weapon | WeaponMaterial
}
