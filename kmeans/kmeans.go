package kmeans

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

// KMeans structure for book keeping the KMeans computation
type KMeans struct {
	maxIterations int
	trainingSet   [][]float64
	labels        []int
	Centroids     [][]float64
	Output        io.Writer
}

const tol = 5

// NewKMeans create a KMeans object to store the result of KMeans clustering
func NewKMeans(k, maxIterations int, trainingSet [][]float64) *KMeans {

	// start all guesses with the zero vector.
	// they will be changed during learning
	guesses := make([]int, len(trainingSet))
	numFeatures := len(trainingSet[0])

	rand.Seed(time.Now().UTC().Unix())
	centroids := make([][]float64, k)
	for i := range centroids {
		centroids[i] = make([]float64, numFeatures)
		copy(centroids[i], trainingSet[rand.Intn(len(trainingSet))])

	}

	return &KMeans{
		maxIterations: maxIterations,

		trainingSet: trainingSet,
		labels:      guesses,

		Centroids: centroids,
		Output:    os.Stdout,
	}

}

// a function to compute the distance between the centroids
func diff(u, v []float64) float64 {
	sum := 0.0
	for i := range u {
		sum += (u[i] - v[i]) * (u[i] - v[i])
	}
	return sum
}

// compute the total distance between previous iteration's cluster centroids
// and current iteration's cluster centroids
func shift(prevcenters, newcenteres [][]float64) float64 {

	var totaldiff float64

	for idx := range prevcenters {
		totaldiff += diff(prevcenters[idx], newcenteres[idx])

	}

	return totaldiff
}

// Labels : Return the cluster labels
func (km *KMeans) Labels() []int {
	return km.labels
}

// Cluster Performs KMeans clustering
func (km *KMeans) Cluster() error {

	if len(km.trainingSet) == 0 {
		err := fmt.Errorf("training data not supplied")
		fmt.Fprint(km.Output, err.Error())
		return err

	}

	numSamples := len(km.trainingSet)
	numClusters := len(km.Centroids)
	numFeatures := len(km.trainingSet[0])

	fmt.Fprintf(km.Output, "Clustering:\n\tModel: K-Means Classification\n\tNumber of Samples: %v\n\tnumFeatures: %v\n\tClasses: %v\n...\n\n",
		numSamples,
		numFeatures,
		numClusters)

	iter := 0
	classCount := make([]int, numClusters)
	classTotal := make([][]float64, numClusters)

	for idx := range classTotal {
		classTotal[idx] = make([]float64, numFeatures)
	}

	prevcenters := make([][]float64, numClusters)

	for idx := range prevcenters {
		prevcenters[idx] = make([]float64, numFeatures)
	}

	for ; iter < km.maxIterations; iter++ {

		for idx := range classTotal {
			classCount[idx] = 0
			for jdx := range classTotal[idx] {
				classTotal[idx][jdx] = 0.0
			}
		}

		for i, x := range km.trainingSet {
			km.labels[i] = 0
			minDiff := diff(x, km.Centroids[0])
			for j := 1; j < numClusters; j++ {
				dist := diff(x, km.Centroids[j])
				if dist < minDiff {
					minDiff = dist
					km.labels[i] = j
				}
			}
			classCount[km.labels[i]]++

			for j := range x {

				classTotal[km.labels[i]][j] += x[j]
			}
		}

		for i := range prevcenters {

			copy(prevcenters[i], km.Centroids[i])

		}

		for j := range km.Centroids {
			if classCount[j] == 0 {
				fmt.Fprintf(km.Output, "Encoutered zero count for cluster=%d\n", j)
				copy(km.Centroids[j], km.trainingSet[rand.Intn(numSamples)])
				continue
			}

			for l := range km.Centroids[j] {

				km.Centroids[j][l] = classTotal[j][l] / float64(classCount[j])

			}

		}

		calctolerance := shift(prevcenters, km.Centroids)

		fmt.Fprintf(km.Output, "iter=%d  shift=%8f\n", iter+1, calctolerance)

		if calctolerance <= tol {

			break
		}

	}

	fmt.Fprintf(km.Output, "clustering finished in iterations:%d\n", iter+1)
	return nil

}
