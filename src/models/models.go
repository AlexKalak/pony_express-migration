package models

type CountryCode struct {
	// gorm.Model
	ID   int    `gorm:"type:BIGINT NOT NULL" json:"-"`
	Code string `gorm:"type:VARCHAR(3) NOT NULL" json:"code"`
}

type Country struct {
	ID            int         `gorm:"type:BIGINT" json:"-"`
	Name          string      `gorm:"type:VARCHAR(255) NOT NULL" json:"name"`
	TrName        string      `gorm:"type:VARCHAR(255) NOT NULL" json:"tr-name"`
	CountryCodeID int         ``
	CountryCode   CountryCode `json:"country-code"`
	RegionID      int         `gorm:"type:BIGINT" json:"region-id"`
	Cities        []City      `json:"cities"`
}

type Area struct {
	ID        int     `gorm:"type:BIGINT" json:"-"`
	Name      string  `gorm:"type:VARCHAR(255) NOT NULL" json:"name" validate:"required"`
	TrName    string  `gorm:"type:VARCHAR(255) NOT NULL" json:"tr-name" validate:"required"`
	Country   Country `json:"country"`
	CountryID int     `gorm:"type:BIGINT" json:"-"`
}

type District struct {
	ID     int    `gorm:"type:BIGINT" json:"id"`
	Name   string `gorm:"type:VARCHAR(255)" json:"name"`
	TrName string `gorm:"type:VARCHAR(255)" json:"tr-name"`
	AreaID int    `gorm:"type:BIGINT" json:"-"`
	Area   Area   `json:"area"`
}

type City struct {
	ID         int      `gorm:"type:BIGINT" json:"-"`
	Name       string   `gorm:"type:VARCHAR(255) NOT NULL" json:"name" validate:"required"`
	TrName     string   `gorm:"type:VARCHAR(255) NOT NULL" json:"tr-name" validate:"required"`
	Region     Region   `json:"region"`
	RegionID   int      `json:"-"`
	Country    Country  `json:"country"`
	CountryID  int      `gorm:"type:BIGINT" json:"-"`
	DistrictID *int     `gorm:"type:BIGINT" json:"-"`
	District   District `json:"district"`
	AreaID     *int     `gorm:"type:BIGINT" json:"-"`
	Area       Area     `json:"area"`
}

type SenderCityWithOffice struct {
	ID     int
	Name   string
	TrName string
}

type SenderCity struct {
	ID                     int    `gorm:"type:BIGINT" json:"-"`
	Name                   string `gorm:"type:VARCHAR(255) NOT NULL"`
	TrName                 string `gorm:"type:VARCHAR(255) NOT NULL"`
	Code                   int
	HasOffice              bool
	SenderRegionID         *int
	SenderRegion           SenderRegion
	SenderCityWithOfficeID int
	SenderCityWithOffice   SenderCityWithOffice
}

type SenderRegion struct {
	ID           int    `gorm:"type:BIGINT"`
	Name         string `gorm:"type:VARCHAR(255) NOT NULL"`
	TrName       string `gorm:"type:VARCHAR(255) NOT NULL"`
	PriceForDoor int    `gorm:"type:BIGINT"`
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
	ID                     int                  `json:"-"`
	WeightID               int                  `json:"-"`
	Weight                 Weight               `json:"weight"`
	RegionID               int                  `json:"-"`
	Region                 Region               `json:"region"`
	PackageTypeID          int                  `json:"-"`
	PackageType            PackageType          `json:"package-type"`
	SenderCityWithOfficeID int                  `json:"-"`
	SenderCityWithOffice   SenderCityWithOffice `json:"sender-city"`
	Price                  int                  `json:"price"`
}

type PriceOverMaxWeight struct {
	ID                     int                  `json:"-"`
	WeightID               int                  `json:"-"`
	Weight                 Weight               `json:"weight"`
	PackageTypeID          int                  `json:"-"`
	PackageType            PackageType          `json:"package-type"`
	RegionID               int                  `json:"-"`
	Region                 Region               `json:"region"`
	SenderCityWithOfficeID int                  `json:"-"`
	SenderCityWithOffice   SenderCityWithOffice `json:"sender-city"`
	Price                  int                  `json:"price"`
}

type ProductType struct {
	// gorm.Model
	ID int `json:"id"`

	EnName   string  `gorm:"type:VARCHAR(255) NOT NULL" json:"en-name" validate:"required,min=2,max=30"`
	RuName   string  `gorm:"type:VARCHAR(255) NOT NULL" json:"ro-name" validate:"required,min=2,max=30"`
	TrName   string  `gorm:"type:VARCHAR(255) NOT NULL" json:"tr-name" validate:"required,min=2,max=30"`
	GtipCode int     `gorm:"type:BIGINT NOT NULL" json:"gtip-code" validate:"required,min=4,numeric"`
	ItemCode int     `gorm:"type:BIGINT NOT NULL" json:"item-code" validate:"required,numeric"`
	Weight   float64 `gorm:"type:DECIMAL(5,2) NOT NULL" json:"weight" validate:"lte=1000,gt=0"`
	Warning  bool    `gorm:"type:bool NOT NULL" json:"warning" validate:"boolean"`
}
