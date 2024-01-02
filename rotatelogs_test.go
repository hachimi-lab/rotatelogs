package rotatelogs

import (
	"log"
	"testing"
	"time"
)

func TestRotateLog(t *testing.T) {
	writer, err := New(
		"./logs/app.log",
		WithMaxAge(time.Minute),
		WithRotateTime(time.Minute),
	)
	if err != nil {
		t.Fatal(err)
	}

	log.SetOutput(writer)

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		log.Println(time.Now().Format(time.DateTime))
	}
}
