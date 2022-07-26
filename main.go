package main

import (
	"log"
	"os"
	"time"

	"github.com/MrDjeb/vk/pkg/config"
	"github.com/MrDjeb/vk/pkg/database"
	"github.com/MrDjeb/vk/pkg/vk"

	"github.com/SevereCloud/vksdk/api"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("[ERROR] ")

	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

	cfg, err := config.Init()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := database.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Photos.DB.Close()

	vkAPI := api.NewVK(cfg.Token)

	vkUPD := vk.NewUPDATer(vkAPI, cfg, db)
	if err := vkUPD.StartWallEditing(); err != nil {
		log.Fatalln(err)
	}
}
