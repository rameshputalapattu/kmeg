package main

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "kmegutil",
		Description: "utility to convert image to kmeg (KMeans expert group) format and reconstruct it as Kmeans color quantized image",
		Usage:       "compress and decompress images using K-means color quantization",
		Commands: []*cli.Command{
			{
				Name:    "compress",
				Aliases: []string{"c"},
				Usage:   "convert an image to kmeg format",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Aliases:  []string{"i"},
						Usage:    "input image file path",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "output kmeg file path",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "levels",
						Aliases:  []string{"l"},
						Usage:    "number of quantization levels",
						Required: true,
						Value:    16,
					},
				},
				Action: compressAction,
			},
			{
				Name:    "decompress",
				Aliases: []string{"d"},
				Usage:   "convert a kmeg file back to an image",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Aliases:  []string{"i"},
						Usage:    "input kmeg file path",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "output image file path",
						Required: true,
					},
				},
				Action: decompressAction,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalf("Error: %v", err)
	}

	logrus.Info("executed the command successfully")
}

func compressAction(c *cli.Context) error {
	input := c.String("input")
	output := c.String("output")
	levels := c.Int("levels")

	if levels <= 0 {
		return errors.New("levels must be greater than 0")
	}

	logrus.Infof("Compressing %s to %s with %d quantization levels", input, output, levels)
	return compress(input, output, levels)
}

func decompressAction(c *cli.Context) error {
	input := c.String("input")
	output := c.String("output")

	logrus.Infof("Decompressing %s to %s", input, output)
	return decompress(input, output)
}
