package models

// CountryInfo holds basic country data used to populate the dashboard.
type CountryInfo struct {
	Name       string
	ISOCode    string
	Capital    string
	Latitude   float64
	Longitude  float64
	Population int
	Area       float64
	Currency   string
}
