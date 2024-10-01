package georm

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkbcommon"
)

func ExampleGeometry_String() {
	point := Geometry[geom.T]{
		Geom: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{42, 42}),
	}

	polygon := Geometry[geom.T]{
		Geom: geom.NewPolygon(geom.XY).MustSetCoords(
			[][]geom.Coord{{{42, 42}, {1, 1}, {2, 2}, {42, 42}}}),
	}

	fmt.Println(point.String())
	fmt.Println(polygon.String())
	// Output: POINT (42 42)
	// POLYGON ((42 42, 1 1, 2 2, 42 42))
}

func TestGeometryStringExpectCannotMarshal(t *testing.T) {
	this := Geometry[geom.T]{nil}
	expect := "cannot marshal geometry: <nil>"
	require.Equal(t, expect, this.String())
}

func TestGeometryValue(t *testing.T) {
	type (
		Input struct {
			Geom geom.T
		}

		Output struct {
			Value driver.Value
			Error error
		}
	)

	tests := []struct {
		Name   string
		Input  Input
		Expect Output
	}{
		{
			Name:   "point (42 42)",
			Input:  Input{Geom: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{42, 42}).SetSRID(4326)},
			Expect: Output{Value: "0101000020e610000000000000000045400000000000004540"},
		},
		{
			Name:   "point (1 2)",
			Input:  Input{Geom: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}).SetSRID(4326)},
			Expect: Output{Value: "0101000020e6100000000000000000f03f0000000000000040"},
		},
		{
			Name:   "expect error unsupported layout",
			Input:  Input{Geom: &geom.Point{}},
			Expect: Output{Error: geom.ErrUnsupportedLayout(geom.NoLayout)},
		},
		{
			Name:   "expect nil nil (value null)",
			Input:  Input{Geom: nil},
			Expect: Output{Value: nil, Error: nil},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var actual Output
			actual.Value, actual.Error = New(test.Input.Geom).Value()

			if test.Expect.Error != nil {
				assert.ErrorIs(t, actual.Error, test.Expect.Error)
				return
			} else {
				require.NoError(t, actual.Error)
			}

			assert.Equal(t, test.Expect.Value, actual.Value)
		})
	}
}

func TestGeometryScan(t *testing.T) {
	type (
		Input struct {
			Geom  geom.T
			value interface{}
		}

		Output struct {
			Geom  geom.T
			Error error
		}
	)

	tests := []struct {
		Name   string
		Input  Input
		Expect Output
	}{
		{
			Name: "expect point (42 42)",
			Input: Input{
				Geom:  &geom.Point{},
				value: "0101000020e610000000000000000045400000000000004540",
			},
			Expect: Output{
				Geom: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{42, 42}).SetSRID(4326),
			},
		},
		{
			Name: "expect point (1 2)",
			Input: Input{
				Geom:  &geom.Point{},
				value: "0101000020e6100000000000000000f03f0000000000000040",
			},
			Expect: Output{
				Geom: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1, 2}).SetSRID(4326),
			},
		},
		{
			Name: "expect point from byte slice (42 42)",
			Input: Input{
				Geom:  &geom.Point{},
				value: []byte{1, 1, 0, 0, 32, 230, 16, 0, 0, 0, 0, 0, 0, 0, 0, 69, 64, 0, 0, 0, 0, 0, 0, 69, 64},
			},
			Expect: Output{
				Geom: geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{42, 42}).SetSRID(4326),
			},
		},
		{
			Name:   "expect err unsupported geometry type",
			Input:  Input{Geom: &geom.Point{}, value: uint(42)},
			Expect: Output{Error: ErrUnexpectedGeometryType},
		},
		{
			Name:   "expect err on decode hex",
			Input:  Input{Geom: &geom.Point{}, value: "notHexString"},
			Expect: Output{Error: hex.InvalidByteError('n')},
		},
		{
			Name:   "expect err on decode wkb",
			Input:  Input{Geom: &geom.Point{}, value: "000000000000"},
			Expect: Output{Error: wkbcommon.ErrUnsupportedType(0)},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var actual Output
			actualThis := New(test.Input.Geom)

			actual.Error = actualThis.Scan(test.Input.value)

			if test.Expect.Error != nil {
				assert.ErrorIs(t, actual.Error, test.Expect.Error)
				return
			} else {
				require.NoError(t, actual.Error)
			}

			assert.Equal(t, test.Expect.Geom, actualThis.Geom)
		})
	}
}

func TestGeometryScanExpectUnexpectedGeometryType(t *testing.T) {
	this := New(&geom.Polygon{})

	err := this.Scan("0101000020e6100000000000000000f03f0000000000000040")
	require.ErrorIs(t, err, ErrUnexpectedValueType)
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
		{Geom: nil, Expect: "geometry"}, // any geometry
	}

	for _, test := range tests {
		t.Run(test.Expect, func(t *testing.T) {
			assert.Equal(t, test.Expect, New(test.Geom).GormDataType())
		})
	}
}
