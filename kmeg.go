package kmeg

import (
	"compress/gzip"
	"encoding/binary"
	"io"

	"github.com/rameshputalapattu/kmeg/kmeans"
)

// Kmeg Struct to hold the Kmeg information
type Kmeg struct {
	Quantlevels [][]uint8
	Labels      []int
	Dx          int
	Dy          int
}

// convert the data type of cluster centroids into uint8
func makeCenters(centroids [][]float32) [][]uint8 {

	quantLevels := len(centroids)

	centers := make([][]uint8, quantLevels)

	for idx := range centers {
		var features []uint8
		for _, feature := range centroids[idx] {
			features = append(features, uint8(feature))
		}
		centers[idx] = features

	}

	return centers

}

// MakeTrainingSet create a training set (samples) from the image pixel array
func MakeTrainingSet(data [][]uint8) [][]float32 {

	trainingSet := make([][]float32, len(data))

	for idx := range trainingSet {
		var features []float32
		for _, feature := range data[idx] {
			features = append(features, float32(feature))
		}

		trainingSet[idx] = features

	}

	return trainingSet

}

// Deflate Take in the kmeans object (Which has image information as the training samples)
// and return Kmeg representation
func Deflate(km *kmeans.KMeans, dx int) (*Kmeg, error) {

	err := km.Cluster()

	if err != nil {
		return nil, err
	}

	kmeg := Kmeg{
		Quantlevels: makeCenters(km.Centroids),
		Labels:      km.Labels(),
		Dx:          dx,
		Dy:          len(km.Labels()) / dx,
	}

	return &kmeg, nil

}

// Inflate returns RGBImage representation from the kmeg representation
func Inflate(kmeg *Kmeg) *RGBImage {

	dx := kmeg.Dx

	dy := len(kmeg.Labels) / dx

	Pix := make([][]uint8, len(kmeg.Labels))

	for idx := range Pix {

		Pix[idx] = kmeg.Quantlevels[kmeg.Labels[idx]]

	}

	rgb := RGBImage{
		Pix: Pix,
		Dx:  dx,
		Dy:  dy,
	}
	return &rgb

}

// Encode : Serialize kmeg representation into kmeg binary format
func Encode(w io.Writer, kmg *Kmeg) error {

	gzw := gzip.NewWriter(w)

	defer gzw.Close()

	//gzw := w

	err := binary.Write(gzw, binary.BigEndian, int64(kmg.Dx))

	if err != nil {
		return err
	}

	err = binary.Write(gzw, binary.BigEndian, int64(kmg.Dy))

	if err != nil {
		return err
	}

	for _, label := range kmg.Labels {

		err := binary.Write(gzw, binary.BigEndian, uint8(label))
		if err != nil {

			return err

		}
	}

	for _, quant := range kmg.Quantlevels {
		for _, color := range quant {

			err := binary.Write(gzw, binary.BigEndian, uint8(color))

			if err != nil {
				return err
			}

		}

	}

	return nil

}

// Decode : Deserialize the kmeg binary format to in-memory Kmeg structure
func Decode(r io.Reader) (*Kmeg, error) {
	var dx int64

	var err error

	gzr, err := gzip.NewReader(r)
	defer gzr.Close()

	if err != nil {
		return nil, err
	}

	err = binary.Read(gzr, binary.BigEndian, &dx)

	if err != nil {
		return nil, err
	}

	var dy int64

	err = binary.Read(gzr, binary.BigEndian, &dy)

	if err != nil {
		return nil, err
	}

	labels := make([]uint8, int(dx)*int(dy))

	err = binary.Read(gzr, binary.BigEndian, labels)

	if err != nil {
		return nil, err
	}

	var quantizedLevels [][]uint8

	for err != io.EOF {

		var pix []uint8

		for i := 0; i < 3; i++ {
			var color uint8
			err = binary.Read(gzr, binary.BigEndian, &color)
			if err == io.EOF {
				break
			}

			if err != nil && err != io.EOF {
				return nil, err
			}

			pix = append(pix, color)

		}

		if err == io.EOF {
			continue
		}

		quantizedLevels = append(quantizedLevels, pix)

	}

	var labelints []int

	for _, elem := range labels {
		labelints = append(labelints, int(elem))
	}

	kmg := Kmeg{
		Quantlevels: quantizedLevels,
		Labels:      labelints,
		Dx:          int(dx),
		Dy:          int(dy),
	}

	return &kmg, nil

}
