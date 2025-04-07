//package models

//import (
//	"time"
//)

//type CountryInfo struct {
//	Area             int                `firestore:"area" json:"area"`
//	Capital          string             `firestore:"capital" json:"capital"`
//	Coordinates      Coordinates        `firestore:"coordinates" json:"coordinates"`
//	LastRetrieval    time.Time          `firestore:"last_retrieval" json:"last_retrieval"`
//	Name             string             `firestore:"name" json:"name"`
//	Population       int                `firestore:"population" json:"population"`
//	Precipitation    float64            `firestore:"precipitation" json:"precipitation"`
//	TargetCurrencies map[string]float64 `firestore:"target_currencies" json:"target_currencies"`
//	Temperature      float64            `firestore:"temperature" json:"temperature"`
//	ISOCode          string             `firestore:"iso_code" json:"iso_code"`
//}

//type Coordinates struct {
//	Latitude  float64 `firestore:"latitude" json:"latitude"`
//	Longitude float64 `firestore:"longitude" json:"longitude"`
//}
