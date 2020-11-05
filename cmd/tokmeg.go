package main

import (
	"context"
	"errors"
	"flag"
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

//TokmegCommand Command to convert images from png format to kmeg binary format
type TokmegCommand struct {
	Params *CmdParams
}

const toKmegHelp = `convert a image into kmeg format`

//Name Gives the command Name
func (cmd *TokmegCommand) Name() string { return "toKmeg" }

//Args returns the command args
func (cmd *TokmegCommand) Args() string { return "" }

//ShortHelp returns short help text
func (cmd *TokmegCommand) ShortHelp() string { return toKmegHelp }

//LongHelp returns long help text
func (cmd *TokmegCommand) LongHelp() string { return toKmegHelp }

//Hidden returns whether it is a hidden command
func (cmd *TokmegCommand) Hidden() bool { return false }

//Register Registers the flag set
func (cmd *TokmegCommand) Register(fs *flag.FlagSet) {

}

//Run Run the ToKmeg command to convert a image into kmeg binary format
func (cmd *TokmegCommand) Run(ctx context.Context, args []string) error {

	if len(cmd.Params.SrcImageFile) == 0 {
		return errors.New("Source Image file should be provided for tokmeg")
	}

	if len(cmd.Params.DstImageFile) == 0 {
		return errors.New("Destination Image file should be provided for tokmeg")
	}

	if cmd.Params.QuantLevels == 0 {
		return errors.New("levels should not be zero")
	}

	return compress(cmd.Params.SrcImageFile, cmd.Params.DstImageFile, cmd.Params.QuantLevels)

}

func readImage(r io.Reader, extn string) (image.Image, error) {

	switch extn {
	case "png":
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
		return ew.Wrapf(err, "error opening the original image file %s\n", imgName)
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

	err = kmeg.Encode(w, kmegImg)

	if err != nil {
		return ew.Wrap(err, "encoding kmeg to disk failed")
	}

	return nil

}
