package geo

import (
	"fmt"
	"math"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
)

const (
	earthRadiusKm = 6371.0
)

// GeoEncode encodes a latitude and longitude into a hex string
func GeoEncode(lat, long float64) (string, error) {
	g := geom.NewPoint(geom.XY).MustSetCoords([]float64{long, lat}).SetSRID(4326)
	return ewkbhex.Encode(g, ewkbhex.NDR)

}

// GeoDecode decodes an array of bytes into a latitude and longitude
func GeoDecodeBytes(location []byte) (*geom.Point, error) {
	return GeoDecodeString(string(location))
}

// GeoDecodeString decodes a hex string into a latitude and longitude
func GeoDecodeString(location string) (*geom.Point, error) {
	locationGeom, err := ewkbhex.Decode(location)
	if err != nil {
		return nil, err
	}

	locationPoint, ok := locationGeom.(*geom.Point)
	if !ok {
		return nil, fmt.Errorf("location is not a point")
	}
	return locationPoint, nil

}

// CalculateDistance calculates the distance between two points specified by latitude and longitude in kilometers
func CalculateDistance(lat1, lon1, lat2, lon2 float64) int {
	rad := math.Pi / 180
	dLat := (lat2 - lat1) * rad
	dLon := (lon2 - lon1) * rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return int(earthRadiusKm * c)
}
