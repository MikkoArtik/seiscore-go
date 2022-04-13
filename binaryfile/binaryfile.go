package binaryfile


import (
	"math"
	"time"
	"os"
	"strings"
	"fmt"
	"strconv"
	"path"
	"errors"
)




type FileHeader struct {
	channelsCount uint16
	frequency uint16
	datetimeStart time.Time
	longitude float64
	latitude float64
}


func formatDuration(days int, hours int, minutes int, seconds int) string {
	duration := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	if days > 0 {
		duration = fmt.Sprintf("%d days ", days) + duration
	}
	return duration
}


type FileInfo struct {
	path string
	formatType string
	frequency int64
	timeStart time.Time
	timeStop time.Time
	longitude float64
	latitude float64
}

func (fileInfo FileInfo) name() string {
	return path.Base(fileInfo.path)
}

func (fileInfo FileInfo) secondsDuration() float64 {
	return fileInfo.timeStop.Sub(fileInfo.timeStart).Seconds()
}

func (fileInfo FileInfo) formattedDuration() string {
	diff := fileInfo.timeStop.Sub(fileInfo.timeStart)
	daysDiff := int(diff.Hours() / 24)
	hoursDiff := int(int(diff.Hours()) % 24)
	minutesDiff := int(int(diff.Minutes()) % 60)
	secondsDiff := int (int(diff.Seconds()) % 60)
	return formatDuration(daysDiff, hoursDiff, minutesDiff, secondsDiff)
}


const BAIKAL7_FMT, BAIKAL8_FMT, SIGMA_FMT = "Baikal7", "Baikal8", "Sigma"
const BAIKAL7_EXTENSION, BAIKAL8_EXTENSION, SIGMA_EXTENSION = "00", "xx", "bin"

var BINARY_FILE_FORMATS = map[string]string{
	BAIKAL7_FMT: BAIKAL7_EXTENSION,
	BAIKAL8_FMT: BAIKAL8_EXTENSION,
	SIGMA_FMT: SIGMA_EXTENSION,
}

const SIGMA_SECONDS_OFFSET = 2
const COMPONENTS_ORDER = "ZXY"


func truncate(num float64, precision float64) float64 {
	return math.Round(num * math.Pow(10, precision)) / math.Pow(10, precision)
}


func isBinaryFilePath(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	splitPath := strings.Split(path, ".")
	extension := splitPath[len(splitPath) - 1]

	allowedExtensions := [3]string{BAIKAL7_EXTENSION, BAIKAL8_EXTENSION, SIGMA_EXTENSION}
	for _, item := range allowedExtensions {
		if item == extension {
			return true
		}
	}
	return false
}


// func readBinary(file *os.File, bytesCount int, skippingBytes int) []byte {
// 	bytes := make([]byte, bytesCount)
// 	file.ReadAt(bytes, int64(skippingBytes))
// 	return bytes
// }


func getDatetimeStartBaikal7(timeBegin uint64) time.Time {
	seconds := timeBegin / 256000000
	nanoseconds := int(float64(timeBegin % 256000000) * math.Pow(10, 9))

	datetimeStart := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	datetimeStart = datetimeStart.Add(time.Second * time.Duration(seconds))
	datetimeStart = datetimeStart.Add(time.Nanosecond * time.Duration(nanoseconds))
	return datetimeStart
}


// type CharType struct {
// 	filePointer *os.File
// 	skippingBytes int
// 	elementsCount int
// }

// func (dataType CharType) byteSize() int {
// 	return 1
// }

// func (dataType CharType) convert() string {
// 	bytesVal := readBinary(
// 		dataType.filePointer, dataType.byteSize() * dataType.elementsCount, 
// 		dataType.skippingBytes)
// 	return string(bytesVal)
// }


// type UnsignedShortType struct {
// 	filePointer *os.File
// 	skippingBytes int
// 	elementsCount int	
// }

// func (dataType UnsignedShortType) byteSize() int {
// 	return 2
// }

// func (dataType UnsignedShortType) getBytes() []byte {
// 	return readBinary(
// 		dataType.filePointer, dataType.byteSize() * dataType.elementsCount, 
// 		dataType.skippingBytes)
// }

// func (dataType UnsignedShortType) convertToNumber() uint16 {
// 	buffer := bytes.NewBuffer(dataType.getBytes())
// 	var result uint16
// 	binary.Read(buffer, binary.LittleEndian, &result)
// 	return result
// }

// func (dataType UnsignedShortType) convertToArray() []uint16 {
// 	buffer := bytes.NewBuffer(dataType.getBytes())
// 	result := make([]uint16, dataType.elementsCount)
// 	binary.Read(buffer, binary.LittleEndian, &result)
// 	return result
// }


// type UnsignedIntType struct {
// 	filePointer *os.File
// 	skippingBytes int
// 	elementsCount int
// }

// func (dataType UnsignedIntType) byteSize() int {
// 	return 4
// }

// func (dataType UnsignedIntType) getBytes() []byte {
// 	return readBinary(
// 		dataType.filePointer, dataType.byteSize() * dataType.elementsCount, 
// 		dataType.skippingBytes)
// }

// func (dataType UnsignedIntType) convertToNumber() int32 {
// 	buffer := bytes.NewBuffer(dataType.getBytes())
// 	var result int32
// 	binary.Read(buffer, binary.LittleEndian, &result)
// 	return result
// }

// func (dataType UnsignedIntType) convertToArray() []int32 {
// 	buffer := bytes.NewBuffer(dataType.getBytes())
// 	result := make([]int32, dataType.elementsCount)
// 	binary.Read(buffer, binary.LittleEndian, &result)
// 	return result
// }


// type DoubleType struct {
// 	filePointer *os.File
// 	skippingBytes int
// 	elementsCount int
// }

// func (dataType DoubleType) byteSize() int {
// 	return 8
// }

// func (dataType DoubleType) getBytes() []byte {
// 	return readBinary(
// 		dataType.filePointer, dataType.byteSize() * dataType.elementsCount, 
// 		dataType.skippingBytes)
// }

// func (dataType DoubleType) convertToNumber() float64 {
// 	buffer := bytes.NewBuffer(dataType.getBytes())
// 	var result float64
// 	binary.Read(buffer, binary.LittleEndian, &result)
// 	return result
// }

// func (dataType DoubleType) convertToArray() []float64 {
// 	buffer := bytes.NewBuffer(dataType.getBytes())
// 	result := make([]float64, dataType.elementsCount)
// 	binary.Read(buffer, binary.LittleEndian, &result)
// 	return result
// }


// type LongType struct {
// 	filePointer *os.File
// 	skippingBytes int
// 	elementsCount int
// }

// func (dataType LongType) byteSize() int {
// 	return 8
// }

// func (dataType LongType) getBytes() []byte {
// 	return readBinary(
// 		dataType.filePointer, dataType.byteSize() * dataType.elementsCount, 
// 		dataType.skippingBytes)
// }

// func (dataType LongType) convertToNumber() uint64 {
// 	buffer := bytes.NewBuffer(dataType.getBytes())
// 	var result uint64
// 	binary.Read(buffer, binary.LittleEndian, &result)
// 	return result
// }

// func (dataType LongType) convertToArray() []uint64 {
// 	buffer := bytes.NewBuffer(dataType.getBytes())
// 	result := make([]uint64, dataType.elementsCount)
// 	binary.Read(buffer, binary.LittleEndian, &result)
// 	return result
// }


func ReadBaikal7Header(path string) FileHeader {
	file, _ := os.Open(path)
	defer file.Close()

	channelsCount := UnsignedShortType{file, 0, 1}.convertToNumber()
	frequency := UnsignedShortType{file, 22, 1}.convertToNumber()

	srcCoords := DoubleType{file, 72, 2}.convertToArray()
	longitude, latitude := truncate(srcCoords[1], 5), truncate(srcCoords[0], 5)
	timeBegin := LongType{file, 104, 1}.convertToNumber()
	datetimeStart := getDatetimeStartBaikal7(timeBegin)
	return FileHeader{channelsCount, frequency, datetimeStart, longitude, latitude} 
}


func ReadBaikal8Header(path string) FileHeader {
	file, _ := os.Open(path)
	defer file.Close()

	channelsCount := UnsignedShortType{file, 0, 1}.convertToNumber()
	dateSrc := UnsignedShortType{file, 6, 3}.convertToArray()

	srcVals := DoubleType{file, 48, 2}.convertToArray()
	frequency := uint16(1 / srcVals[0])
	seconds := int(srcVals[1])
	nanoseconds := int((srcVals[1] - float64(seconds)) * math.Pow(10, 9))

	datetimeStart := time.Date(
		int(dateSrc[2]), time.Month(dateSrc[1]), int(dateSrc[0]), 
		0, 0, 0, 0, time.UTC)
	datetimeStart = datetimeStart.Add(time.Second * time.Duration(seconds))
	datetimeStart = datetimeStart.Add(time.Nanosecond * time.Duration(nanoseconds))
	
	srcCoords := DoubleType{file, 72, 2}.convertToArray()
	longitude, latitude := truncate(srcCoords[0], 5), truncate(srcCoords[1], 5)
	return FileHeader{channelsCount, frequency, datetimeStart, longitude, latitude} 
}


func ReadSigmaHeader(path string) FileHeader {
	file, _ := os.Open(path)
	defer file.Close()

	channelsCount := UnsignedShortType{file, 12, 1}.convertToNumber()
	frequency := UnsignedShortType{file, 24, 1}.convertToNumber()
	latitudeSrc, longitudeSrc := CharType{file, 40, 8}.convert(), CharType{file, 48, 9}.convert()

	datetimeSrc := UnsignedIntType{file, 60, 2}.convertToArray()
	dateLine := strconv.FormatInt(int64(datetimeSrc[0]), 10)
	timeLine := fmt.Sprintf("%06d", int64(datetimeSrc[1]))

	year, _ := strconv.ParseInt(dateLine[:2], 10, 64)
	year += 2000

	month, _ := strconv.ParseInt(dateLine[2:4], 10, 64)
	day, _ := strconv.ParseInt(dateLine[4:], 10, 64)
	hours, _ := strconv.ParseInt(timeLine[:2], 10, 64)
	minutes, _ := strconv.ParseInt(timeLine[2:4], 10, 64)
	seconds, _ := strconv.ParseInt(timeLine[4:], 10, 64)

	datetimeStart := time.Date(
		int(year), time.Month(month), int(day), 
		int(hours), int(minutes), int(seconds), 0, time.UTC)

	
	integerPart, _ := strconv.ParseFloat(longitudeSrc[:3], 64)
	decimalPart, _ := strconv.ParseFloat(longitudeSrc[3:len(longitudeSrc) - 1], 64)
	longitude := truncate(integerPart + decimalPart / 60, 5)

	integerPart, _ = strconv.ParseFloat(latitudeSrc[:2], 64)
	decimalPart, _ = strconv.ParseFloat(latitudeSrc[2:len(latitudeSrc) - 1], 64)
	latitude := truncate(integerPart + decimalPart / 60, 5)
	return FileHeader{channelsCount, frequency, datetimeStart, longitude, latitude}
}


func resampling(signal []int, resample_parameter int) []int {
	resample_arr_size := len(signal) / resample_parameter
	var result []int
	for i := 0; i < resample_arr_size; i++ {
		sum_val := 0
		for j := 0; j < resample_parameter; j++ {
			index := i * resample_parameter + j
			sum_val += signal[index]
		}
		value := sum_val / resample_parameter
		result = append(result, value)
	}
	return result
}


type BinaryFile struct {
	path string
	resampleFrequency int
	isUseAvgValues bool
}

func (binFile BinaryFile) fileExtention() (string, error) {
	if len(binFile.path) == 0 {
		return "", errors.New("Empty file path")
	}
	splitPath := strings.Split(binFile.path, ".")
	return splitPath[len(splitPath) - 1], nil
}
