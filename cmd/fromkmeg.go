package main

import (
	"context"
	"errors"
	"flag"
	"image/png"
	"os"

	ew "github.com/pkg/errors"
	"github.com/rameshputalapattu/kmeg"
)

//FromkmegCommand Command to convert images from png format to kmeg binary format
type FromkmegCommand struct {
	Params *CmdParams
}

const fromKmegHelp = `convert a image into kmeg format`

//Name Gives the command Name
func (cmd *FromkmegCommand) Name() string { return "fromKmeg" }

//Args returns the command args
func (cmd *FromkmegCommand) Args() string { return "" }

//ShortHelp returns short help text
func (cmd *FromkmegCommand) ShortHelp() string { return fromKmegHelp }

//LongHelp returns long help text
func (cmd *FromkmegCommand) LongHelp() string { return fromKmegHelp }

//Hidden returns whether it is a hidden command
func (cmd *FromkmegCommand) Hidden() bool { return false }

//Register Registers the flag set
func (cmd *FromkmegCommand) Register(fs *flag.FlagSet) {

}

//Run Run the ToKmeg command to convert a image into kmeg binary format
func (cmd *FromkmegCommand) Run(ctx context.Context, args []string) error {

	if len(cmd.Params.SrcImageFile) == 0 {
		return errors.New("Source Image file should be provided for tokmeg")
	}

	if len(cmd.Params.DstImageFile) == 0 {
		return errors.New("Destination Image file should be provided for tokmeg")
	}

	return decompress(cmd.Params.SrcImageFile, cmd.Params.DstImageFile)

}

func decompress(decompName string, reconstName string) error {

	r, err := os.Open(decompName)
	if err != nil {
		return ew.Wrap(err, "opening the .kmeg file for reading failed")
	}

	w, err := os.Create(reconstName)

	if err != nil {
		return ew.Wrap(err, "creating the file to write re-constructed image failed")
	}

	defer w.Close()

	kmegImg, err := kmeg.Decode(r)

	if err != nil {
		return ew.Wrap(err, "Decoding kmeg to rgb image failed")
	}

	rgb := kmeg.Inflate(kmegImg)

	rgba := kmeg.ConvertToRGBA(rgb)

	err = png.Encode(w, rgba)

	if err != nil {
		return ew.Wrap(err, "encoding to png format failed")
	}

	return nil

}
