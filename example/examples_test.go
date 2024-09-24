package example

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-geom"

	"github.com/ybru-tech/georm"
)

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

func TestExamplePoint(t *testing.T) {
	tx := db.Debug()

	err := tx.AutoMigrate(
		// CREATE TABLE "addresses" ("id" bigserial,"address" text,"geo_point" Geometry(Point, 4326),PRIMARY KEY ("id"))
		Address{},
		// CREATE TABLE "zones" ("id" bigserial,"title" text,"geo_polygon" Geometry(Polygon, 4326),PRIMARY KEY ("id"))
		Zone{},
	)
	require.NoError(t, err)

	// INSERT INTO "addresses" ("address","geo_point") VALUES ('some address','010100000000000000000045400000000000003840') RETURNING "id"
	err = tx.Create(&Address{
		Address: "some address",
		GeoPoint: georm.Point{
			Geom: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{42, 24}),
		},
	}).Error
	require.NoError(t, err)

	err = tx.Create(&Zone{
		Title: "some zone",
		GeoPolygon: georm.Polygon{
			Geom: geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{
				{{30, 10}, {40, 40}, {20, 40}, {10, 20}, {30, 10}},
			}),
		},
	}).Error
	// INSERT INTO "zones" ("title","geo_polygon") VALUES ('some zone','010300000001000000050000000000000000003e4000000000000024400000000000004440000000000000444000000000000034400000000000004440000000000000244000000000000034400000000000003e400000000000002440') RETURNING "id"
	require.NoError(t, err)

	point := georm.Point{
		Geom: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{25, 26}).SetSRID(georm.SRID),
	}
	var result Zone

	// SELECT * FROM "zones" WHERE ST_Contains(geo_polygon, '0101000020e610000000000000000039400000000000003a40') ORDER BY "zones"."id" LIMIT 1
	err = tx.Model(&Zone{}).
		Where("ST_Contains(geo_polygon, ?)", point).
		First(&result).Error
	require.NoError(t, err)
	fmt.Printf("%#v\n", result)
}

func TestStringer(t *testing.T) {
	// POINT (25 26)
	fmt.Println(georm.Point{
		Geom: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{25, 26}).SetSRID(georm.SRID),
	})

	// POLYGON ((30 10, 40 40, 20 40, 10 20, 30 10))
	fmt.Println(georm.Polygon{
		Geom: geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{
			{{30, 10}, {40, 40}, {20, 40}, {10, 20}, {30, 10}},
		}),
	})
}
