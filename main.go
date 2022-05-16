package main

import (
	"fmt"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"math"
)

var earthRadius = 6.371 * math.Pow(10, 6)

var locationMap = make(map[string][]Location, 0)

type Location struct {
	UserId int
	Name   string
	Point  Point
}

type Point struct {
	Lat float64
	Lng float64
}

func NewPoint(lat, lng float64) Point {
	return Point{
		Lat: lat,
		Lng: lng,
	}
}

func main() {
	SaveLocation(1, NewPoint(40.390411191986914, 49.80271891796518), "insaatcilar")
	SaveLocation(2, NewPoint(40.389381278047914, 49.80419480691158), "bizim market")
	SaveLocation(3, NewPoint(40.37933891912101, 49.804308324196), "asan 3")
	SaveLocation(4, NewPoint(40.37489145979754, 49.815389147652944), "elmler m")
	SaveLocation(5, NewPoint(40.36601912643649, 49.831608592207175), "iceriseher")
	result := Search(NewPoint(40.37927587970779, 49.829811861062765), 5000)
	fmt.Println(result)
}

const level = 13

func SaveLocation(userId int, point Point, name string) {
	latLng := s2.LatLngFromDegrees(point.Lat, point.Lng)
	cellID := s2.CellIDFromLatLng(latLng)
	cellIDOnStorageLevel := cellID.Parent(level).String()
	locationMap[cellIDOnStorageLevel] = append(locationMap[cellIDOnStorageLevel], Location{
		UserId: userId,
		Point:  point,
		Name:   name,
	})
}

func Search(point Point, radius uint32) []Location {
	latLng := s2.LatLngFromDegrees(point.Lat, point.Lng)
	centerPoint := s2.PointFromLatLng(latLng)
	centerAngle := float64(radius) / earthRadius
	capFromCenterAngle := s2.CapFromCenterAngle(centerPoint, s1.Angle(centerAngle))
	rc := s2.RegionCoverer{MaxLevel: level, MinLevel: level}
	cu := rc.Covering(capFromCenterAngle)
	var result []Location
	counterEx := 0
	for _, cellId := range cu {
		if cellId.Level() < level {
			counter := 0
			for ci := cellId.ChildBeginAtLevel(level); ci < cellId.ChildEndAtLevel(level).Next(); ci = ci.Next() {
				result = appendLocation(ci.String(), result)
				counter++
			}
			fmt.Println(counter)
		} else {
			result = appendLocation(cellId.String(), result)
			counterEx++
		}
	}
	fmt.Println("ex counter ", counterEx)
	return result
}

func appendLocation(cellId string, result []Location) []Location {
	if locations, ok := locationMap[cellId]; ok {
		for _, location := range locations {
			result = append(result, location)
		}
	}
	return result
}
