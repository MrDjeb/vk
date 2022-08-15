package vk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MrDjeb/vk/pkg/database"
	"github.com/SevereCloud/vksdk/api"
)

type ResponseURL struct {
	Server      int    `json:"server"`
	Photo       string `json:"photo"`
	Mid         int    `json:"mid"`
	Hash        string `json:"hash"`
	MessageCode int    `json:"message_code"`
	ProfileAid  int    `json:"profile_aid"`
	PhotoList   string `json:"photos_list"`
}

func PostUpload(urlPath, photoPath, field string) (*ResponseURL, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile(field, photoPath)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(photoPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(fw, file)
	if err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", urlPath, bytes.NewReader(body.Bytes()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		logVK.Printf("PostUpload failed with response code: %d\n", res.StatusCode)
	}

	Rbody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var data = new(ResponseURL)
	err = json.Unmarshal(Rbody, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func MakePhotoAlbum(photoPath string, AlbumID int) error {

	url, err := vkUPD.API.PhotosGetUploadServer(api.Params{
		"album_id": AlbumID,
	})
	if err != nil {
		return err
	}
	logVK.Printf("Get URL: %s\n", url.UploadURL)

	rUrl, err := PostUpload(url.UploadURL, photoPath, "file1")
	if err != nil {
		return err
	}
	logVK.Printf("Request.server after POST: %d\n", rUrl.Server)

	rSave, err := vkUPD.API.PhotosSave(api.Params{
		"album_id":    AlbumID,
		"server":      rUrl.Server,
		"photos_list": rUrl.PhotoList,
		"hash":        rUrl.Hash,
		"latitude":    "90",
		"longitude":   "180",
		"caption":     "STONER",
	})
	if err != nil {
		return err
	}
	if err = vkUPD.DB.Photos.Insert(database.Photo{PhotoID: rSave[0].ID}); err != nil {
		return err
	}
	logVK.Printf("Photo ID after save: %d. Uploud to album(%d) succses\n", rSave[0].ID, AlbumID)

	time.Sleep(201 * time.Millisecond)

	return nil
}

func MakePhotoWall(photoPath string) error {

	url, err := vkUPD.API.PhotosGetWallUploadServer(api.Params{
		"group_id": vkUPD.CFG.MyID,
	})
	if err != nil {
		return err
	}
	logVK.Printf("Get URL: %s\n", url.UploadURL)

	rUrl, err := PostUpload(url.UploadURL, photoPath, "photo")
	if err != nil {
		return err
	}
	logVK.Printf("Request.server after POST: %d\n", rUrl.Server)

	rSave, err := vkUPD.API.PhotosSaveWallPhoto(api.Params{
		"user_id":  vkUPD.CFG.MyID,
		"group_id": vkUPD.CFG.MyID,
		"server":   rUrl.Server,
		"photo":    rUrl.Photo,
		"hash":     rUrl.Hash,
	})
	if err != nil {
		return err
	}
	logVK.Printf("Photo ID after save: %d\n", rSave[0].ID)

	rPost, err := vkUPD.API.WallPost(api.Params{
		"owner_id":           vkUPD.CFG.MyID,
		"friends_only":       0,
		"attachments":        "photo" + strconv.Itoa(vkUPD.CFG.MyID) + "_" + strconv.Itoa(rSave[0].ID),
		"mute_notifications": 1,
		"copyright":          "https://kinggizzardandthelizardwizard.com/releases",
	})
	if err != nil {
		return err
	}
	if err = vkUPD.DB.Photos.Insert(database.Photo{PhotoID: rSave[0].ID}); err != nil {
		return err
	}
	logVK.Printf("MakePhotoWall succses, postID: %d", rPost.PostID)

	time.Sleep(500 * time.Millisecond)

	return nil
}

func MakePost(PhotoID int) (int, error) {
	res, err := vkUPD.API.WallPost(api.Params{
		"owner_id":           vkUPD.CFG.MyID,
		"friends_only":       0,
		"attachments":        "photo" + strconv.Itoa(vkUPD.CFG.MyID) + "_" + strconv.Itoa(PhotoID),
		"mute_notifications": 1,
		"copyright":          "https://kinggizzardandthelizardwizard.com/releases",
		"message":            time.Now().Format("Mon Jan 2 15:04:05"),
	})
	if err != nil {
		return 0, err
	}
	logVK.Printf("MakePost succses, postID: %d", res.PostID)

	time.Sleep(500 * time.Millisecond)

	return res.PostID, nil
}

func PinPost(PostID int) error {
	res, err := vkUPD.API.WallPin(api.Params{
		"owner_id": vkUPD.CFG.MyID,
		"post_id":  PostID,
	})
	if err != nil {
		return err
	}
	if res != 1 {
		return errors.New("wallPin return response different from 1, can't WallPin")
	}
	logVK.Printf("PinPost succses, postID: %d", PostID)

	time.Sleep(500 * time.Millisecond)

	return nil
}

func EditPost(PhotoID, PostID int) error {
	rPost, err := vkUPD.API.WallEdit(api.Params{
		"owner_id":           vkUPD.CFG.MyID,
		"post_id":            PostID,
		"friends_only":       0,
		"attachments":        "photo" + strconv.Itoa(vkUPD.CFG.MyID) + "_" + strconv.Itoa(PhotoID),
		"mute_notifications": 1,
		"copyright":          "https://kinggizzardandthelizardwizard.com/releases",
		"message":            time.Now().Format("Mon Jan 2 15:04:05"),
	})
	if err != nil {
		return err
	}

	logVK.Printf("Edit succses, postID: %d", rPost.PostID)

	time.Sleep(500 * time.Millisecond)

	return nil
}

func RemovePhoto(PhotoID int) error {
	del, err := vkUPD.API.PhotosDelete(api.Params{
		"owner_id": vkUPD.CFG.MyID,
		"photo_id": PhotoID,
	})
	if err != nil {
		return err
	}
	if del != 1 {
		logVK.Printf("Exept removing last photo by ID: %d\n", PhotoID)
	} else {
		vkUPD.DB.Photos.Delete(PhotoID)
		logVK.Printf("Remove succses: %d\n", PhotoID)
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}

func ClearPhotos() error {
	photos, err := vkUPD.DB.Photos.ReadAll()
	if err != nil {
		return err
	}
	for _, photo := range photos {
		if err := RemovePhoto(photo.PhotoID); err != nil {
			return err
		}
	}
	return nil
}

func SetOnline() error {
	res, err := vkUPD.API.AccountSetOnline(api.Params{
		"voip": 0,
	})
	if err != nil {
		return err
	}
	if res != 1 {
		return errors.New("setOnline return response different from 1, can't setOnline")
	}
	logVK.Println("SetOnline succses")

	time.Sleep(500 * time.Millisecond)

	return nil
}
