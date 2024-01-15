package rotatelogs_test

import (
	"github.com/hachimi-lab/rotatelogs"

	"log"
	"testing"
	"time"
)

func TestRotateLog(t *testing.T) {
	writer, err := rotatelogs.New(
		"./logs/app.log",
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithTimePeriod(rotatelogs.Daily),
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
