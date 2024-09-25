package examples

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-geom"
	"gorm.io/gorm"

	"github.com/ybru-tech/georm"
)

func equalTableWithGeometries(t *testing.T, expect, actual TableWithAllGeometries) {
	t.Helper()

	assert.Equal(t, expect.Point, actual.Point)
	assert.Equal(t, expect.LineString, actual.LineString)
	assert.Equal(t, expect.Polygon, actual.Polygon)
	assert.Equal(t, expect.MultiPoint, actual.MultiPoint)
	assert.Equal(t, expect.MultiLineString, actual.MultiLineString)
	assert.Equal(t, expect.MultiPolygon, actual.MultiPolygon)
	assert.Equal(t, expect.GeometryCollection, actual.GeometryCollection)
}

type TableWithAllGeometries struct {
	gorm.Model

	Point              georm.Point
	LineString         georm.LineString
	Polygon            georm.Polygon
	MultiPoint         georm.MultiPoint
	MultiLineString    georm.MultiLineString
	MultiPolygon       georm.MultiPolygon
	GeometryCollection georm.GeometryCollection
}

func TestCRUDTableWithAllGeometries(t *testing.T) {
	var (
		// test data
		_point1 = georm.Point{Geom: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{42, 42}).SetSRID(4326)}
		_point2 = georm.Point{Geom: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 1}).SetSRID(4326)}

		_lineString1 = georm.LineString{Geom: geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{42, 42}, {1, 1}}).SetSRID(4326)}
		_lineString2 = georm.LineString{Geom: geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{1, 1}, {2, 2}}).SetSRID(4326)}

		_polygon1 = georm.Polygon{Geom: geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{{{42, 42}, {1, 1}, {2, 2}, {42, 42}}}).SetSRID(4326)}
		_polygon2 = georm.Polygon{Geom: geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{{{1, 1}, {2, 2}, {3, 3}, {1, 1}}}).SetSRID(4326)}

		_multiPoint1 = georm.MultiPoint{Geom: geom.NewMultiPoint(geom.XY).MustSetCoords([]geom.Coord{{42, 42}, {1, 1}}).SetSRID(4326)}
		_multiPoint2 = georm.MultiPoint{Geom: geom.NewMultiPoint(geom.XY).MustSetCoords([]geom.Coord{{1, 1}, {2, 2}}).SetSRID(4326)}

		_multiLineString1 = georm.MultiLineString{Geom: geom.NewMultiLineString(geom.XY).MustSetCoords([][]geom.Coord{{{42, 42}, {1, 1}, {2, 2}, {42, 42}}}).SetSRID(4326)}
		_multiLineString2 = georm.MultiLineString{Geom: geom.NewMultiLineString(geom.XY).MustSetCoords([][]geom.Coord{{{1, 1}, {2, 2}, {3, 3}, {1, 1}}}).SetSRID(4326)}

		_multiPolygon1 = georm.MultiPolygon{Geom: geom.NewMultiPolygon(geom.XY).MustSetCoords([][][]geom.Coord{{{{42, 42}, {1, 1}, {2, 2}, {42, 42}}}}).SetSRID(4326)}
		_multiPolygon2 = georm.MultiPolygon{Geom: geom.NewMultiPolygon(geom.XY).MustSetCoords([][][]geom.Coord{{{{1, 1}, {2, 2}, {3, 3}, {1, 1}}}}).SetSRID(4326)}

		_geometryCollection1 = georm.GeometryCollection{Geom: geom.NewGeometryCollection().MustPush(
			geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{42, 42}),
			geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{42, 42}, {1, 1}}),
		).SetSRID(4326)}
		_geometryCollection2 = georm.GeometryCollection{Geom: geom.NewGeometryCollection().MustPush(
			geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 1}),
			geom.NewLineString(geom.XY).MustSetCoords([]geom.Coord{{1, 1}, {2, 2}}),
		).SetSRID(4326)}
	)

	migrator := db.Migrator()

	err := migrator.AutoMigrate(&TableWithAllGeometries{})
	require.NoError(t, err)

	var (
		object TableWithAllGeometries

		objectForCreate = TableWithAllGeometries{
			Point:              _point1,
			LineString:         _lineString1,
			Polygon:            _polygon1,
			MultiPoint:         _multiPoint1,
			MultiLineString:    _multiLineString1,
			MultiPolygon:       _multiPolygon1,
			GeometryCollection: _geometryCollection1,
		}

		objectForUpdate = TableWithAllGeometries{
			Point:              _point2,
			LineString:         _lineString2,
			Polygon:            _polygon2,
			MultiPoint:         _multiPoint2,
			MultiLineString:    _multiLineString2,
			MultiPolygon:       _multiPolygon2,
			GeometryCollection: _geometryCollection2,
		}
	)

	// Create geometries
	err = db.Create(&objectForCreate).Error
	require.NoError(t, err)

	// Get created geometries
	err = db.First(&object, objectForCreate.ID).Error
	require.NoError(t, err)

	equalTableWithGeometries(t, objectForCreate, object)

	// Update geometries
	err = db.Model(&object).Updates(objectForUpdate).Error
	require.NoError(t, err)

	// Get updated geometries
	err = db.First(&object, objectForCreate.ID).Error
	require.NoError(t, err)

	equalTableWithGeometries(t, objectForUpdate, object)

	// Delete geometries
	err = db.Delete(&object).Error
	require.NoError(t, err)

	// Try get deleted geometries
	err = db.First(&object, objectForCreate.ID).Error
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)

	err = migrator.DropTable(&TableWithAllGeometries{})
	require.NoError(t, err)
}
