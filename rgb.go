package kmeg

import (
	"image"
	"image/draw"
)

//RGBImage :holds data for pixels of image in RGB format
//Pixels is in matrix form  - each row holds color levels of a single pixel
type RGBImage struct {
	Pix [][]uint8
	Dx  int
	Dy  int
}

//ConvertToRGB :converts from RGBA (Red Green Blue Alpha ) format to RGB (Red Green Blue Format)
//Drops Alpha value from the color vector (color vector is of length 3  for RGB vs 4 in RGBA)
func ConvertToRGB(img image.Image) *RGBImage {

	switch src := img.(type) {
	case *image.NRGBA64, *image.NRGBA, *image.RGBA64, *image.NYCbCrA, *image.CMYK, *image.YCbCr:
		b := src.Bounds()
		imgRgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(imgRgba, imgRgba.Bounds(), src, b.Min, draw.Src)
		return convertToRGB(imgRgba)

	case *image.RGBA:
		return convertToRGB(src)

	default:
		return nil

	}

}

func convertToRGB(img *image.RGBA) *RGBImage {

	Dx, Dy := img.Rect.Dx(), img.Rect.Dy()

	Pix := make([][]uint8, Dx*Dy)

	for idx := range Pix {

		colors := img.Pix[idx*4 : idx*4+4]
		rgbcolors := colors[:3]
		Pix[idx] = rgbcolors

	}

	rgb := RGBImage{
		Pix: Pix,
		Dx:  Dx,
		Dy:  Dy,
	}

	return &rgb

}

//ConvertToRGBA converts from RGB (Red Green Blue ) format to RGB (Red Green Blue Alpha Format)
//In RGBA format,all Pixels are stored in a single sequence. Each slice of 4 consecutive
//elements represent a single pixel RGBA values. A constant alpha level of 255 is appended
//to each pixel's RGB value
func ConvertToRGBA(rgb *RGBImage) *image.RGBA {

	Dx, Dy := rgb.Dx, rgb.Dy

	Rect := image.Rect(0, 0, Dx, Dy)

	Stride := rgb.Dx * 4

	//Pix := make([]uint8,Dx*Dy*4)
	var Pix []uint8

	for _, pix := range rgb.Pix {

		rgbacolors := append(pix, uint8(255))

		Pix = append(Pix, rgbacolors...)

	}

	rgba := image.RGBA{
		Pix:    Pix,
		Stride: Stride,
		Rect:   Rect,
	}

	return &rgba

}
