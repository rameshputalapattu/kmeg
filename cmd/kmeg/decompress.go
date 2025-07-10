package main

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	ew "github.com/pkg/errors"
	"github.com/rameshputalapattu/kmeg"
)

func writeImage(img *image.RGBA, w io.Writer, extn string) error {
	switch extn {
	case ".png":
		err := png.Encode(w, img)
		return err
	case ".jpg", ".JPEG", ".jpeg", ".JPG":
		err := jpeg.Encode(w, img, &jpeg.Options{Quality: 100})
		return err
	default:
		return errors.New("not a valid extension")
	}
}

func decompress(decompName string, reconstName string) error {
	r, err := os.Open(decompName)
	if err != nil {
		return ew.Wrap(err, "opening the .kmeg file for reading failed")
	}
	defer r.Close()

	w, err := os.Create(reconstName)
	if err != nil {
		return ew.Wrap(err, "creating the file to write re-constructed image failed")
	}
	defer w.Close()

	kmegImg, err := kmeg.Decode(r)
	if err != nil {
		return ew.Wrap(err, "decoding kmeg to rgb image failed")
	}

	rgb := kmeg.Inflate(kmegImg)
	rgba := kmeg.ConvertToRGBA(rgb)
	ext := filepath.Ext(reconstName)

	err = writeImage(rgba, w, ext)
	if err != nil {
		return ew.Wrap(err, "encoding to image format failed")
	}

	return nil
}
