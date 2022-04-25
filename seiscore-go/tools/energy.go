package tools


import (
	"math"
	"gonum.org/v1/gonum/integrate"
)


type Limit struct {
	Low float64
	High float64
}


func GetAmplitudeEnergy(signal []int32, normingCoeff float64) float64 {
	var energy float64
	for i := 0; i < len(signal); i++ {
		discrete := float64(signal[i]) / normingCoeff
		energy += math.Pow(discrete, normingCoeff)
	}
	return energy
}


func GetSpectrumEnergy(spectrum [][]float64, frequencyLimit Limit) float64 {
	if frequencyLimit.Low == 0 && frequencyLimit.High == 0 {
		return 0
	}

	xs, ys := []float64{}, []float64{}
	for i := 0; i < len(spectrum); i++ {
		frequency := spectrum[i][0]
		if frequency < frequencyLimit.Low {
			continue
		}
		if frequency > frequencyLimit.High {
			break	
		}
		xs = append(xs, frequency)

		amplitude := spectrum[i][1]
		ys = append(ys, amplitude)
	}
	return integrate.Trapezoidal(xs, ys)
}