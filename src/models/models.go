package models

type CountryCode struct {
	// gorm.Model
	ID   int    `gorm:"type:BIGINT NOT NULL" json:"-"`
	Code string `gorm:"type:VARCHAR(3) NOT NULL" json:"code"`
}

type Country struct {
	ID            int         `gorm:"type:BIGINT" json:"-"`
	Name          string      `gorm:"type:VARCHAR(255) NOT NULL" json:"name"`
	CountryCodeID int         ``
	CountryCode   CountryCode `json:"country-code"`
	RegionID      int         `gorm:"type:BIGINT" json:"region-id"`
	Cities        []City      `json:"cities"`
}

type Area struct {
	ID        int     `gorm:"type:BIGINT" json:"-"`
	Name      string  `gorm:"type:VARCHAR(255) NOT NULL" json:"name" validate:"required"`
	Country   Country `json:"country"`
	CountryID int     `gorm:"type:BIGINT" json:"-"`
}

type District struct {
	ID     int    `gorm:"type:BIGINT" json:"id"`
	Name   string `gorm:"type:VARCHAR(255)" json:"name"`
	AreaID int    `gorm:"type:BIGINT" json:"-"`
	Area   Area   `json:"area"`
}

type City struct {
	ID         int      `gorm:"type:BIGINT" json:"-"`
	Name       string   `gorm:"type:VARCHAR(255) NOT NULL" json:"name" validate:"required"`
	Region     Region   `json:"region"`
	RegionID   int      `json:"-"`
	Country    Country  `json:"country"`
	CountryID  int      `gorm:"type:BIGINT" json:"-"`
	DistrictID *int     `gorm:"type:BIGINT" json:"-"`
	District   District `json:"district"`
	AreaID     *int     `gorm:"type:BIGINT" json:"-"`
	Area       Area     `json:"area"`
}

type SenderCity struct {
	ID   int    `gorm:"type:BIGINT" json:"-"`
	Name string `gorm:"type:VARCHAR(255) NOT NULL" json:"name" validate:"required"`
}

type SenderRegion struct {
	ID           int        `gorm:"type:BIGINT" json:"-"`
	Name         string     `gorm:"type:VARCHAR(255) NOT NULL" json:"name" validate:"required"`
	SenderCityID int        `grom:"type:BIGINT" json:"-"`
	SenderCity   SenderCity `json:"sender-city"`
	PriceForDoor int        `gorm:"type:BIGINT" json:"price-for-door"`
}

type Region struct {
	ID   int    `gorm:"type:BIGINT" json:"id"`
	Name string `gorm:"type:VARCHAR(255)" json:"name"`
}

type Weight struct {
	ID     int     `json:""`
	Weight float64 `json:"weight"`
}

type PackageType struct {
	ID   int    `json:"-"`
	Name string `json:"name"`
}

type DeliveryType struct {
	ID   int    `gorm:"type:BIGINT" json:"-"`
	Name string `json:"name" validate:"required,min=3"`
}

type Price struct {
	ID            int         `json:"-"`
	WeightID      int         `json:"-"`
	Weight        Weight      `json:"weight"`
	RegionID      int         `json:"-"`
	Region        Region      `json:"region"`
	PackageTypeID int         `json:"-"`
	PackageType   PackageType `json:"package-type"`
	SenderCityID  int         `json:"-"`
	SenderCity    SenderCity  `json:"sender-city"`
	Price         int         `json:"price"`
}

type PriceOverMaxWeight struct {
	ID            int         `json:"-"`
	WeightID      int         `json:"-"`
	Weight        Weight      `json:"weight"`
	PackageTypeID int         `json:"-"`
	PackageType   PackageType `json:"package-type"`
	RegionID      int         `json:"-"`
	Region        Region      `json:"region"`
	SenderCityID  int         `json:"-"`
	SenderCity    SenderCity  `json:"sender-city"`
	Price         int         `json:"price"`
}
