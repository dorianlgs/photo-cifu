package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func CreateGallery(e *core.RequestEvent, app *pocketbase.PocketBase) error {

	name := e.Request.FormValue("name")

	location := e.Request.FormValue("location")

	mf, mh, err := e.Request.FormFile("imagesZip")

	if err != nil {
		return e.BadRequestError("Failed to read zip file", err)
	}

	tmf, tmh, err := e.Request.FormFile("thumbnail")

	if err != nil {
		return e.BadRequestError("Failed to read zip file", err)
	}

	archive, err := zip.NewReader(mf, int64(mh.Size))
	if err != nil {
		return e.BadRequestError("Failed to open the zip file", err)
	}

	images_collection, err := app.FindCollectionByNameOrId("images")
	if err != nil {
		return err
	}

	galleries_collection, err := app.FindCollectionByNameOrId("galleries")
	if err != nil {
		return err
	}

	galleryRecord := core.NewRecord(galleries_collection)

	transactErr := app.RunInTransaction(func(txApp core.App) error {

		var imagesNames []string

		for _, f := range archive.File {

			fileInArchive, err := f.Open()
			if err != nil {
				return fmt.Errorf("Failed to open the zip file "+f.Name, err)
			}

			imgRecord := core.NewRecord(images_collection)

			imgRecord.Set("likes", 0)
			dataFile := new(bytes.Buffer)
			if _, err := io.Copy(dataFile, fileInArchive); err != nil {
				return fmt.Errorf("Failed to copy the file "+f.Name, err)
			}

			fileInArchive.Close()

			f2, _ := filesystem.NewFileFromBytes(dataFile.Bytes(), f.Name)

			imgRecord.Set("image", f2)

			err = txApp.Save(imgRecord)
			if err != nil {
				return fmt.Errorf("Failed to save the image "+f.Name, err)
			}

			imagesNames = append(imagesNames, string(imgRecord.Id))

		}

		galleryRecord.Set("name", name)
		galleryRecord.Set("location", location)

		thumbnailData, err := io.ReadAll(tmf)
		if err != nil {
			return fmt.Errorf("Failed to read thumbnail file", err)
		}

		f3, _ := filesystem.NewFileFromBytes(thumbnailData, tmh.Filename)

		galleryRecord.Set("images", imagesNames)

		galleryRecord.Set("thumbnail", f3)

		err = txApp.Save(galleryRecord)
		if err != nil {
			return fmt.Errorf("Failed to save the gallery "+name, err)
		}

		return nil
	})

	if transactErr != nil {
		return e.BadRequestError("Failed to create the gallery", transactErr)
	}

	return e.JSON(http.StatusOK, map[string]any{"galleryId": galleryRecord.Id})
}
