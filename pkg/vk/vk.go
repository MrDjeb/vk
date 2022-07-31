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

var logVK *log.Logger

//stupid posting: uploud photo from PC---wall post---delete wall post
func (u *UPDATer) StartStupidWallPosting() error {
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
			if err := MakePhotoWall(vkUPD.CFG.DirParh + file.Name()); err != nil {
				return err
			}
			<-t.C
			t.Reset(time.Duration(rand.Intn(vkUPD.CFG.DelMax-vkUPD.CFG.DelMin)+vkUPD.CFG.DelMin+vkUPD.CFG.Delay) * time.Millisecond)
		}
	}
}

//Upload local dir to album
func (u *UPDATer) StartAlbumLoad() error {
	logVK = log.New(os.Stderr, "[VK] ", log.LstdFlags|log.Lmsgprefix)

	files, err := ioutil.ReadDir(vkUPD.CFG.DirParh)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := MakePhotoAlbum(vkUPD.CFG.DirParh+file.Name(), vkUPD.CFG.AlbumID); err != nil {
			return err
		}
	}

	return nil
}

func (u *UPDATer) StartWallEditing() error {
	logVK = log.New(os.Stderr, "[VK] ", log.LstdFlags|log.Lmsgprefix)

	t := time.NewTicker(time.Duration(vkUPD.CFG.Delay) * time.Millisecond)
	defer t.Stop()

	photos, err := vkUPD.DB.Photos.ReadAll()
	if err != nil {
		return err
	}

	for {
		for _, photo := range photos {
			if err := EditPost(photo.PhotoID); err != nil {
				return err
			}
			<-t.C
			t.Reset(time.Duration(rand.Intn(vkUPD.CFG.DelMax-vkUPD.CFG.DelMin)+vkUPD.CFG.DelMin+vkUPD.CFG.Delay) * time.Millisecond)
		}
	}
}
