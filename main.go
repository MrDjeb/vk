package main

import (
	"log"
	"time"

	"github.com/MrDjeb/vk/pkg/config"
	"github.com/MrDjeb/vk/pkg/database"
	"github.com/MrDjeb/vk/pkg/vk"

	"github.com/SevereCloud/vksdk/api"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("[ERROR] ")

	var err error
	time.Local, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Printf("error loading location \"Europe/Moscow\": %v\n", err)
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
