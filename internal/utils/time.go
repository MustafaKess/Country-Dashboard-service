package utils

import (
	"fmt"
	"time"
)

// CustomTime wraps time.Time to allow custom JSON formatting
type CustomTime struct {
	time.Time
}

// Layout with timezone – adjust as needed
const timeLayout = "2006-01-02 15:04:05 MST"

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	// Convert to Europe/Oslo timezone
	loc, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return nil, err
	}
	formatted := ct.Time.In(loc).Format(timeLayout)
	return []byte(fmt.Sprintf(`"%s"`, formatted)), nil
}
