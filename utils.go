package sentimenter

import (
	"encoding/json"
	"log"

	uuid "github.com/satori/go.uuid"
)

func getNewID() string {
	return uuid.NewV4().String()
}

func serializeOrFail(o interface{}) []byte {
	b, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
