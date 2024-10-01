package georm

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/wkt"
)

var (
	ErrUnexpectedGeometryType = errors.New("unexpected geometry type")
	ErrUnexpectedValueType    = errors.New("unexpected value type")
)

var SRID = 4326

type (
	Geometry[T geom.T] struct{ Geom T }

	Point              = Geometry[*geom.Point]
	LineString         = Geometry[*geom.LineString]
	Polygon            = Geometry[*geom.Polygon]
	MultiPoint         = Geometry[*geom.MultiPoint]
	MultiLineString    = Geometry[*geom.MultiLineString]
	MultiPolygon       = Geometry[*geom.MultiPolygon]
	GeometryCollection = Geometry[*geom.GeometryCollection]
)

func New[T geom.T](geom T) Geometry[T] { return Geometry[T]{geom} }

// Scan impl sql.Scanner
func (g *Geometry[T]) Scan(value interface{}) (err error) {
	var (
		wkb []byte
		ok  bool
	)

	switch v := value.(type) {
	case string:
		wkb, err = hex.DecodeString(v)
	case []byte:
		wkb = v
	default:
		return ErrUnexpectedGeometryType
	}

	if err != nil {
		return err
	}

	geometryT, err := ewkb.Unmarshal(wkb)
	if err != nil {
		return err
	}

	g.Geom, ok = geometryT.(T)
	if !ok {
		return ErrUnexpectedValueType
	}

	return
}

// Value impl driver.Valuer
func (g Geometry[T]) Value() (driver.Value, error) {
	if geom.T(g.Geom) == nil {
		return nil, nil
	}

	sb := &bytes.Buffer{}
	if err := ewkb.Write(sb, binary.LittleEndian, g.Geom); err != nil {
		return nil, err
	}

	return hex.EncodeToString(sb.Bytes()), nil
}

// GormDataType impl schema.GormDataTypeInterface
func (g Geometry[T]) GormDataType() string {
	srid := strconv.Itoa(SRID)

	switch any(g.Geom).(type) {
	case *geom.Point:
		return "Geometry(Point, " + srid + ")"
	case *geom.LineString:
		return "Geometry(LineString, " + srid + ")"
	case *geom.Polygon:
		return "Geometry(Polygon, " + srid + ")"
	case *geom.MultiPoint:
		return "Geometry(MultiPoint, " + srid + ")"
	case *geom.MultiLineString:
		return "Geometry(MultiLineString, " + srid + ")"
	case *geom.MultiPolygon:
		return "Geometry(MultiPolygon, " + srid + ")"
	case *geom.GeometryCollection:
		return "Geometry(GeometryCollection, " + srid + ")"
	default:
		return "geometry"
	}
}

// String returns geometry formatted using WKT format
func (g Geometry[T]) String() string {
	if geomWkt, err := wkt.Marshal(g.Geom); err == nil {
		return geomWkt
	}

	return fmt.Sprintf("cannot marshal geometry: %T", g.Geom)
}
