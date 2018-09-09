package sentimenter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	idPrefix = "tid"
)

// getNewID parses Firestore valid IDs (can't start with digits)
func getNewID() string {
	return fmt.Sprintf("%s-%s", idPrefix, uuid.NewV4().String())
}

func serializeOrFail(o interface{}) []byte {
	b, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func getMapValAsString(m map[string]interface{}, k string) (v string, e error) {
	if val, has := m[k]; has {
		item := val.(map[string]interface{})
		if v, h := item["stringValue"]; h {
			return v.(string), nil
		}
	}
	return "", errors.New("Key not found or an invalid string format")
}

func getMapValAsTimestamp(m map[string]interface{}, k string) (t time.Time, e error) {
	if val, has := m[k]; has {
		item := val.(map[string]interface{})
		if v, h := item["timestampValue"]; h {
			sv := v.(string)
			t1, e := time.Parse(time.RFC3339, sv)
			if e != nil {
				return time.Now(), fmt.Errorf("Invalid time format (RFC3339): %s", sv)
			}
			return t1, nil
		}
	}
	return time.Now(), errors.New("Key not found or an invalid time format")
}
