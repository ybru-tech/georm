package ex_storage

import (
	"gorm.io/gorm"

	"github.com/ybru-tech/georm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) MigrationTables() error {
	return s.db.AutoMigrate(
		&Address{},
		&Zone{},
		&Route{},
	)
}
func (s *Storage) DropTables() error {
	return s.db.Migrator().DropTable(
		&Address{},
		&Zone{},
		&Route{},
	)
}

// Addresses

func (s *Storage) AddAddresses(address ...*Address) error {
	return s.db.Create(address).Error
}
func (s *Storage) GetAddress(id uint) (*Address, error) {
	var address Address

	if err := s.db.First(&address, id).Error; err != nil {
		return nil, err
	}

	return &address, nil
}

// FindAddressesInPolygon finds addresses that are inside a polygon
func (s *Storage) FindAddressesInPolygon(polygon georm.Polygon) ([]Address, error) {
	var addresses []Address

	tx := s.db.
		Model(&Address{}).
		Where("ST_Contains(?, geo_point)", polygon)

	if err := tx.Find(&addresses).Error; err != nil {
		return nil, err
	}

	return addresses, nil
}
func (s *Storage) UpdateAddress(address *Address) error {
	return s.db.Updates(address).Error
}
func (s *Storage) DeleteAddress(id uint) error {
	return s.db.Delete(&Address{}, id).Error
}

// Zones

func (s *Storage) AddZone(zone *Zone) error {
	return s.db.Create(zone).Error
}
func (s *Storage) GetZone(id uint) (*Zone, error) {
	var zone Zone

	if err := s.db.First(&zone, id).Error; err != nil {
		return nil, err
	}

	return &zone, nil
}

// FindZonesContainingPoint finds zones that contain a point
func (s *Storage) FindZonesContainingPoint(point georm.Point) ([]Zone, error) {
	var zones []Zone

	tx := s.db.
		Model(&Zone{}).
		Where("ST_Contains(geom, ?)", point)

	if err := tx.Find(&zones).Error; err != nil {
		return nil, err
	}

	return zones, nil
}
func (s *Storage) UpdateZone(zone *Zone) error {
	return s.db.Updates(zone).Error
}
func (s *Storage) DeleteZone(id uint) error {
	return s.db.Delete(&Zone{}, id).Error
}

// Routes

func (s *Storage) AddRoute(route *Route) error {
	return s.db.Create(route).Error
}
func (s *Storage) GetRoute(id uint) (*Route, error) {
	var route Route

	if err := s.db.First(&route, id).Error; err != nil {
		return nil, err
	}

	return &route, nil
}

// FindRoutesInterZone finds routes that intersect with a zone
func (s *Storage) FindRoutesInterZone(zone *Zone) ([]Route, error) {
	var routes []Route

	tx := s.db.
		Model(&Route{}).
		Where("ST_Intersects(geom, ?)", zone.GeoPolygon)

	if err := tx.Find(&routes).Error; err != nil {
		return nil, err
	}

	return routes, nil
}
func (s *Storage) UpdateRoute(route *Route) error {
	return s.db.Updates(route).Error
}
func (s *Storage) DeleteRoute(id uint) error {
	return s.db.Delete(&Route{}, id).Error
}
