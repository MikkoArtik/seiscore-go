package tools

import "fmt"


type InvalidParameter struct {
	message string
}

func (customError InvalidParameter) Error() string {
	return fmt.Sprintf("InvalidParameter: %s", customError.message)
}


type BadSignalData struct {
	message string
}

func (customError BadSignalData) Error() string {
	return fmt.Sprintf("BadSignalData: %s", customError.message)
}


func Marmett(signal []float64, order uint16) ([]float64, error) {
	filteredSignal := signal

	if order % 2 == 0 {
		return filteredSignal, InvalidParameter{"Invalid marmett order"}
	}

	if len(signal) < 5 {
		return filteredSignal, BadSignalData{"Short signal - length less than 5 discretes"}
	}

	for i := uint16(0); i < order; i++ {
		array := []float64{}

		item := (filteredSignal[0] + filteredSignal[1]) / 2
		array = append(array, item)
		for j := 1; j < len(signal) - 1; j++ {
			item = (filteredSignal[j - 1] + filteredSignal[j + 1]) / 4 + filteredSignal[j] / 2
			array = append(array, item)
		}
		item = (filteredSignal[len(filteredSignal) - 1] + filteredSignal[len(filteredSignal) - 2]) / 2
		array = append(array, item)
		filteredSignal = array
	}
	return filteredSignal, nil
}