package model

import (
	"encoding/json"
	"time"
)

type Geospatial struct {
	ID           uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	GadmID       string      `gorm:"<-:create;unique" json:"-"`
	ParentGadmID string      `gorm:"<-" json:"-"`
	Name         string      `gorm:"<-" json:"name"`
	Type         string      `gorm:"<-" json:"type"`
	Level        uint        `gorm:"<-" json:"level"`
	Geometry     string      `gorm:"type:geometry" json:"-"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Child        *Geospatial `gorm:"-:all" json:"-"`
	ChildJSON    *Geospatial `gorm:"-:all" json:"child,omitempty"`
}

func (g *Geospatial) MarshalJSON() ([]byte, error) {
	type Alias Geospatial
	return json.Marshal(&struct {
		*Alias
		ChildJSON interface{} `json:"child,omitempty"`
	}{
		Alias:     (*Alias)(g),
		ChildJSON: g.Child,
	})
}

type GeospatialFilter struct {
	Name        string   `json:"name"`
	Levels      []uint   `json:"levels"`
	Types       []string `json:"types"`
	ExcludedIds []uint   `json:"excludedIds"`
	ParentIds   []uint   `json:"parentIds"`
	Lat         float64  `json:"lat"`
	Lng         float64  `json:"lng"`
	Nested      bool     `json:"nested"`
}

type GeospatialFilterParams struct {
	Name        string            `query:"name" form:"name"`
	Levels      string            `query:"levels" form:"levels"`
	Types       string            `query:"types" form:"types"`
	LatLng      string            `query:"latlng" form:"latlng"`
	ExcludedIds string            `query:"excludedIds" form:"excludedIds"`
	ParentIds   string            `query:"parentIds" form:"parentIds"`
	Limit       uint              `query:"limit" form:"limit"`
	Page        uint              `query:"page" form:"page"`
	Sort        map[string]string `query:"sort" form:"sort"`
	Nested      bool              `query:"nested" form:"nested"`
}
