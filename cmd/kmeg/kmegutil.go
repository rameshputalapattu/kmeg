package main

import (
	"context"
	"errors"
	"flag"

	"github.com/genuinetools/pkg/cli"
	"github.com/sirupsen/logrus"
)

func main() {

	p := cli.NewProgram()
	p.Name = "kmegutil"
	p.Description = "utility to convert image to kmeg (KMeans expert group) format and reconstruct it as Kmeans color quantized image"

	params := &CmdParams{}

	p.Commands = []cli.Command{&TokmegCommand{params}, &FromkmegCommand{params}}

	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.StringVar(&params.SrcImageFile, "from", "", "source image file")
	p.FlagSet.StringVar(&params.DstImageFile, "to", "", "destination image file")
	p.FlagSet.IntVar(&params.QuantLevels, "levels", 0, "number of quantization levels")

	p.Before = func(ctx context.Context) error {

		if len(params.SrcImageFile) == 0 || len(params.DstImageFile) == 0 {
			return errors.New("Both source and destination image file paths must be provided")
		}

		return nil

	}

	p.Run()
	logrus.Info("executed the command successfully")
}
