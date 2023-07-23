package datetime

import "time"

func IsValidRFC3339(dateString string) bool {
	_, err := time.Parse(time.RFC3339, dateString)
	return err == nil
}
