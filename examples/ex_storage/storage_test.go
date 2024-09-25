package ex_storage

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-geom"
	"gorm.io/gorm"

	"github.com/ybru-tech/georm"
)

func TestStorageAddressCRUD(t *testing.T) {
	address := &Address{
		Address:  "some address",
		GeoPoint: georm.New(geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{42, 42}).SetSRID(4326)),
	}

	err := storage.AddAddresses(address)
	require.NoError(t, err)

	getAddress, err := storage.GetAddress(address.ID)
	require.NoError(t, err)

	require.Equal(t, address, getAddress)

	dataForUpdate := &Address{
		ID:       address.ID,
		Address:  "some address updated",
		GeoPoint: georm.New(geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{32, 32}).SetSRID(4326)),
	}

	err = storage.UpdateAddress(dataForUpdate)
	require.NoError(t, err)

	getAddress, err = storage.GetAddress(address.ID)
	require.NoError(t, err)

	require.Equal(t, dataForUpdate, getAddress)

	err = storage.DeleteAddress(address.ID)
	require.NoError(t, err)

	getAddress, err = storage.GetAddress(address.ID)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestStorage_FindAddressesInPolygon(t *testing.T) {
	addresses := []*Address{
		{Address: "address 1", GeoPoint: georm.New(geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{12, 14}).SetSRID(4326))},
		{Address: "address 2", GeoPoint: georm.New(geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{14, 14}).SetSRID(4326))},
		{Address: "address 3", GeoPoint: georm.New(geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{14, 12}).SetSRID(4326))},
		{Address: "address 4", GeoPoint: georm.New(geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{12, 12}).SetSRID(4326))},
		{Address: "address 5", GeoPoint: georm.New(geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{20, 20}).SetSRID(4326))},
		{Address: "address 6", GeoPoint: georm.New(geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{10, 10}).SetSRID(4326))},
	}

	polygon := georm.Polygon{
		Geom: geom.NewPolygon(geom.XY).MustSetCoords(
			[][]geom.Coord{
				{{11, 11}, {11, 15}, {15, 15}, {15, 11}, {11, 11}},
			},
		).SetSRID(4326),
	}

	err := storage.AddAddresses(addresses...)
	require.NoError(t, err)

	addressesInPolygon, err := storage.FindAddressesInPolygon(polygon)
	require.NoError(t, err)

	expectedAddresses := []*Address{addresses[0], addresses[1], addresses[2], addresses[3]}

	require.Len(t, addressesInPolygon, len(expectedAddresses))
	for _, expectAddress := range expectedAddresses {
		var actualAddress *Address

		for _, address := range expectedAddresses {
			if address.ID != expectAddress.ID {
				continue
			}

			actualAddress = address
		}

		require.NotNilf(t, actualAddress, "expected address %#v not found in actual", expectAddress)
		require.Equal(t, expectAddress, actualAddress)
	}
}
