package georm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-geom"
)

func ExamplePolygon_String() {
	fmt.Println(Polygon{
		Geom: geom.NewPolygon(geom.XY).MustSetCoords(
			[][]geom.Coord{
				{{42, 42}, {1, 1}, {2, 2}, {42, 42}},
			}),
	}.String())

	// Output: POLYGON ((42 42, 1 1, 2 2, 42 42))
}

func ExamplePoint_String() {
	fmt.Println(Point{
		Geom: geom.NewPoint(geom.XY).MustSetCoords(
			geom.Coord{42, 42},
		).SetSRID(4326),
	}.String())

	// Output: POINT (42 42)
}

func TestGeometryValue(t *testing.T) {
	tests := []struct {
		Geom   geom.T
		Expect string
	}{
		{
			Geom:   geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{42, 42}).SetSRID(4326),
			Expect: "0101000020e610000000000000000045400000000000004540",
		},
		{
			Geom:   geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}).SetSRID(4326),
			Expect: "0101000020e6100000000000000000f03f0000000000000040",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("value %T", test.Geom), func(t *testing.T) {
			value, err := New(test.Geom).Value()
			require.NoError(t, err)

			assert.Equal(t, test.Expect, value)
		})
	}
}

func TestGeometryScan(t *testing.T) {
	tests := []struct {
		Geom   geom.T
		value  interface{}
		Expect geom.T
	}{
		{
			Geom:   &geom.Point{},
			value:  "0101000020e610000000000000000045400000000000004540",
			Expect: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{42, 42}).SetSRID(4326),
		},
		{
			Geom:   &geom.Point{},
			value:  "0101000020e6100000000000000000f03f0000000000000040",
			Expect: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}).SetSRID(4326),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("scan %T", test.Geom), func(t *testing.T) {
			g := New(test.Geom)

			err := g.Scan(test.value)
			require.NoError(t, err)

			assert.Equal(t, test.Expect, g.Geom)
		})
	}
}

func TestGeometryGormDataType(t *testing.T) {
	tests := []struct {
		Geom   geom.T
		Expect string
	}{
		{Geom: geom.NewPoint(geom.XY), Expect: "Geometry(Point, 4326)"},
		{Geom: geom.NewLineString(geom.XY), Expect: "Geometry(LineString, 4326)"},
		{Geom: geom.NewPolygon(geom.XY), Expect: "Geometry(Polygon, 4326)"},
		{Geom: geom.NewMultiPoint(geom.XY), Expect: "Geometry(MultiPoint, 4326)"},
		{Geom: geom.NewMultiLineString(geom.XY), Expect: "Geometry(MultiLineString, 4326)"},
		{Geom: geom.NewMultiPolygon(geom.XY), Expect: "Geometry(MultiPolygon, 4326)"},
		{Geom: geom.NewGeometryCollection(), Expect: "Geometry(GeometryCollection, 4326)"},
	}

	for _, test := range tests {
		t.Run("GormDataType_"+test.Expect, func(t *testing.T) {
			assert.Equal(t, test.Expect, New(test.Geom).GormDataType())
		})
	}
}
