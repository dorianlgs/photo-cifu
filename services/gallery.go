package services

import (
	"archive/zip"
	"bytes"
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

	var imagesNames []string

	for _, f := range archive.File {

		fileInArchive, err := f.Open()
		if err != nil {
			return e.BadRequestError("Failed to open the zip file "+f.Name, err)
		}

		collection, err := app.FindCollectionByNameOrId("images")
		if err != nil {
			return err
		}

		record := core.NewRecord(collection)

		record.Set("likes", 0)
		dataFile := new(bytes.Buffer)
		if _, err := io.Copy(dataFile, fileInArchive); err != nil {
			return e.BadRequestError("Failed to copy the file "+f.Name, err)
		}

		fileInArchive.Close()

		f2, _ := filesystem.NewFileFromBytes(dataFile.Bytes(), f.Name)

		record.Set("image", f2)

		err = app.Save(record)
		if err != nil {
			return e.BadRequestError("Failed to save the image "+f.Name, err)
		}

		imagesNames = append(imagesNames, string(record.Id))

	}

	collection, err := app.FindCollectionByNameOrId("galleries")
	if err != nil {
		return err
	}

	record1 := core.NewRecord(collection)

	record1.Set("name", name)
	record1.Set("location", location)

	thumbnailData, err := io.ReadAll(tmf)
	if err != nil {
		return e.BadRequestError("Failed to read thumbnail file", err)
	}

	f3, _ := filesystem.NewFileFromBytes(thumbnailData, tmh.Filename)

	record1.Set("images", imagesNames)

	record1.Set("thumbnail", f3)

	err = app.Save(record1)
	if err != nil {
		return e.BadRequestError("Failed to save the gallery "+name, err)
	}

	return e.JSON(http.StatusOK, map[string]any{"galleryId": record1.Id})
}
