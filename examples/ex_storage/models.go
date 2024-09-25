package ex_storage

import "github.com/ybru-tech/georm"

type Address struct {
	ID       uint `gorm:"primaryKey"`
	Address  string
	GeoPoint georm.Point
}

type Zone struct {
	ID         uint `gorm:"primaryKey"`
	Title      string
	GeoPolygon georm.Polygon
}

type Route struct {
	ID       uint `gorm:"primaryKey"`
	Title    string
	GeoRoute georm.LineString
}
