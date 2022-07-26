package main

import (
	"log"

	"github.com/MrDjeb/vk/pkg/config"
	"github.com/MrDjeb/vk/pkg/database"
	"github.com/MrDjeb/vk/pkg/vk"

	"github.com/SevereCloud/vksdk/api"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("[ERROR] ")

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
	if err := vkUPD.Start(); err != nil {
		log.Fatalln(err)
	}
}
