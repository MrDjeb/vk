package vk

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/MrDjeb/vk/pkg/config"
	"github.com/MrDjeb/vk/pkg/database"
	"github.com/SevereCloud/vksdk/api"
)

var vkUPD *UPDATer

type UPDATer struct {
	API *api.VK
	CFG *config.Config
	DB  database.Tables
}

func NewUPDATer(api *api.VK, cfg *config.Config, db database.Tables) *UPDATer {
	vkUPD = &UPDATer{
		API: api,
		CFG: cfg,
		DB:  db,
	}
	return vkUPD
}

const (
	max = 500
	min = -500
)

var logVK *log.Logger

func (u *UPDATer) Start() error {
	logVK = log.New(os.Stderr, "[VK] ", log.LstdFlags|log.Lmsgprefix)

	if err := ClearPhotos(); err != nil {
		return err
	}

	t := time.NewTicker(time.Duration(vkUPD.CFG.Delay) * time.Millisecond)
	defer t.Stop()

	files, err := ioutil.ReadDir(vkUPD.CFG.DirParh)
	if err != nil {
		return err
	}

	for {
		for _, file := range files {
			if err := ClearPhotos(); err != nil {
				return err
			}
			if err := MakePhoto(vkUPD.CFG.DirParh + file.Name()); err != nil {
				return err
			}
			<-t.C
			t.Reset(time.Duration(rand.Intn(max-min)+min+vkUPD.CFG.Delay) * time.Millisecond)
		}
	}
}
