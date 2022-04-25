package binaryfile

import (
	"fmt"
)


type BadHeaderData struct {
	message string
}

func (customError BadHeaderData) Error() string {
	return fmt.Sprintf("BadHeaderData: %s", customError.message)
}


type BadFilePath struct {
	message string
}

func (customError BadFilePath) Error() string {
	return fmt.Sprintf("BadFilePath: %s", customError.message)
} 


type InvalidResampleFrequency struct {
	message string
}

func (customError InvalidResampleFrequency) Error() string {
	return fmt.Sprintf("InvalidResampleFrequency: %s", customError.message)
}


type InvalidDatetimeValue struct {
	message string
}

func (customError InvalidDatetimeValue) Error() string {
	return fmt.Sprintf("InvalidDatetimeValue: %s", customError.message)
}


type UnknownComponentName struct {
	message string
}

func (customError UnknownComponentName) Error() string {
	return fmt.Sprintf("UnknownComponentName: %s", customError.message)
}


type BadSignalData struct {
	message string
}

func (customError BadSignalData) Error() string {
	return fmt.Sprintf("BadSignalData: %s", customError.message)
}