package api

import (
	"github.com/aminkbi/microChatApp/internal/util"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	util.InitLogger()

	err := util.ConnectMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	os.Exit(code)

}
