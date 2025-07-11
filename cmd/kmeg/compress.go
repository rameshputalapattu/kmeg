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
	"github.com/rameshputalapattu/kmeg/kmeans"
)

func readImage(r io.Reader, extn string) (image.Image, error) {
	switch extn {
	case ".png":
		img, err := png.Decode(r)
		return img, err
	case ".jpg", ".JPEG", ".jpeg", ".JPG":
		img, err := jpeg.Decode(r)
		return img, err
	default:
		return nil, errors.New("not a valid extension")
	}
}

func compress(imgName string, compImgName string, quantlevels int) error {
	r, err := os.Open(imgName)
	if err != nil {
		return ew.Wrapf(err, "error opening the original image file %s", imgName)
	}
	defer r.Close()

	ext := filepath.Ext(imgName)
	img, err := readImage(r, ext)
	if err != nil {
		return ew.Wrap(err, "error decoding the image")
	}

	rgb := kmeg.ConvertToRGB(img)
	km := kmeans.NewKMeans(quantlevels, 300, kmeg.MakeTrainingSet(rgb.Pix))
	kmegImg, err := kmeg.Deflate(km, rgb.Dx)
	if err != nil {
		return ew.Wrap(err, "construction of kmeg format failed")
	}

	w, err := os.Create(compImgName)
	if err != nil {
		return ew.Wrap(err, "creating .kmeg file failed")
	}
	defer w.Close()

	err = kmeg.Encode(w, kmegImg)
	if err != nil {
		return ew.Wrap(err, "encoding kmeg to disk failed")
	}

	return nil
}
