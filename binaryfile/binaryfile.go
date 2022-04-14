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


func getDatetimeStartBaikal7(timeBegin uint64) time.Time {
	seconds := timeBegin / 256000000
	nanoseconds := int(float64(timeBegin % 256000000) * math.Pow(10, 9))

	datetimeStart := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	datetimeStart = datetimeStart.Add(time.Second * time.Duration(seconds))
	datetimeStart = datetimeStart.Add(time.Nanosecond * time.Duration(nanoseconds))
	return datetimeStart
}


func getDatetimeStartSigma(dateNum int32, timeNum int32) (time.Time, error) {
	dateLine := strconv.FormatInt(int64(dateNum), 10)
	if len(dateLine) != 6 {
		return time.Time{}, errors.New("Invalid date in header")
	}

	year, _ := strconv.ParseInt(dateLine[:2], 10, 64)
	year += 2000

	month, _ := strconv.ParseInt(dateLine[2:4], 10, 64)
	if month < 1 || month > 12 {
		return time.Time{}, errors.New("Invalid month in header")
	}

	day, _ := strconv.ParseInt(dateLine[4:], 10, 64)

	firstMonthDay := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastMonthDay := firstMonthDay.AddDate(0, 1, 0).Add(-time.Nanosecond)
	if day < 1 || int(day) > lastMonthDay.Day() {
		return time.Time{}, errors.New("Invalid day in header")
	}

	timeLine := fmt.Sprintf("%06d", int64(timeNum))
	hours, _ := strconv.ParseInt(timeLine[:2], 10, 64)
	if hours > 23 {
		return time.Time{}, errors.New("Invalid hours in header")
	}

	minutes, _ := strconv.ParseInt(timeLine[2:4], 10, 64)
	if minutes > 59 {
		return time.Time{}, errors.New("Invalid minutes in header")
	}

	seconds, _ := strconv.ParseInt(timeLine[4:], 10, 64)
	if seconds > 59 {
		return time.Time{}, errors.New("Invalid seconds in header")
	}

	datetimeStart := time.Date(
		int(year), time.Month(month), int(day), 
		int(hours), int(minutes), int(seconds), 0, time.UTC)
	return datetimeStart, nil
}


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


func ReadSigmaHeader(path string) (FileHeader, error) {
	file, _ := os.Open(path)
	defer file.Close()

	channelsCount := UnsignedShortType{file, 12, 1}.convertToNumber()
	frequency := UnsignedShortType{file, 24, 1}.convertToNumber()
	latitudeSrc, longitudeSrc := CharType{file, 40, 8}.convert(), CharType{file, 48, 9}.convert()

	datetimeSrc := UnsignedIntType{file, 60, 2}.convertToArray()
	datetimeStart, err := getDatetimeStartSigma(datetimeSrc[0], datetimeSrc[1])
	if err != nil {
		return FileHeader{}, err
	}
	
	integerPart, _ := strconv.ParseFloat(longitudeSrc[:3], 64)
	decimalPart, _ := strconv.ParseFloat(longitudeSrc[3:len(longitudeSrc) - 1], 64)
	longitude := truncate(integerPart + decimalPart / 60, 5)

	integerPart, _ = strconv.ParseFloat(latitudeSrc[:2], 64)
	decimalPart, _ = strconv.ParseFloat(latitudeSrc[2:len(latitudeSrc) - 1], 64)
	latitude := truncate(integerPart + decimalPart / 60, 5)

	return FileHeader{channelsCount, frequency, datetimeStart, longitude, latitude}, nil
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
	Path string
	ResampleFrequency uint16
	IsUseAvgValues bool
}

func (binFile BinaryFile) fileExtension() (string, error) {
	if len(binFile.Path) == 0 {
		return "", errors.New("Empty file path")
	}
	splitPath := strings.Split(binFile.Path, ".")
	return splitPath[len(splitPath) - 1], nil
}

func (binFile BinaryFile) formatType() (string, error) {
	currentExtension, err := binFile.fileExtension()
	if err != nil {
		return "", err
	}

	for fileFormat, fileExtension := range BINARY_FILE_FORMATS {
		if currentExtension == fileExtension {
			return fileFormat, nil
		}
	}
	return "", errors.New("Invalid file format")
}

func (binFile BinaryFile) fileHeader() (FileHeader, error) {
	formatType, err := binFile.formatType()
	if err != nil {
		return FileHeader{}, err
	}

	switch formatType {
	case BAIKAL7_FMT:
		return ReadBaikal7Header(binFile.Path), nil
	case BAIKAL8_FMT:
		return ReadBaikal8Header(binFile.Path), nil
	case SIGMA_FMT:
		header, err := ReadSigmaHeader(binFile.Path)
		if err != nil {
			return FileHeader{}, err
		}
		return header, nil
	default:
		return FileHeader{}, errors.New("Unknown format type")
	}
}

func (binFile BinaryFile) originFrequency() (uint16, error) {
	fileHeader, err := binFile.fileHeader()
	if err != nil {
		return 0, errors.New("Bad file header format")
	}
	return fileHeader.frequency, nil
}

func (binFile BinaryFile) GetResampleFrequency() (uint16, error) {
	originFrequency, err := binFile.originFrequency()
	if err != nil {
		return 0, err
	}

	switch freq := binFile.ResampleFrequency; {
	case freq < 0:
		return 0, errors.New("Invalid resample frequency value")
	case freq == 0:
		return originFrequency, nil
	case originFrequency % freq == 0:
		return freq, nil
	default:
		return 0, errors.New("Invalid resample frequency value")
	}
}

