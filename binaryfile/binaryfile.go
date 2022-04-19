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


type Coordinate struct {
	Longitude float64
	Latitude float64
}


type FileHeader struct {
	frequency uint16
	datetimeStart time.Time
	coordinate Coordinate
}


func formatDuration(days int, hours int, minutes int, seconds int) string {
	duration := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	if days > 0 {
		duration = fmt.Sprintf("%d days ", days) + duration
	}
	return duration
}


type FileInfo struct {
	Path string
	FormatType string
	Frequency uint16
	TimeStart time.Time
	TimeStop time.Time
	Coordinate Coordinate
}

func (fileInfo FileInfo) name() string {
	return path.Base(fileInfo.Path)
}

func (fileInfo FileInfo) secondsDuration() float64 {
	return fileInfo.TimeStop.Sub(fileInfo.TimeStart).Seconds()
}

func (fileInfo FileInfo) formattedDuration() string {
	diff := fileInfo.TimeStop.Sub(fileInfo.TimeStart)
	daysDiff := int(diff.Hours() / 24)
	hoursDiff := int(int(diff.Hours()) % 24)
	minutesDiff := int(int(diff.Minutes()) % 60)
	secondsDiff := int (int(diff.Seconds()) % 60)
	return formatDuration(daysDiff, hoursDiff, minutesDiff, secondsDiff)
}


const (
	BAIKAL7_FMT, BAIKAL8_FMT, SIGMA_FMT = "Baikal7", "Baikal8", "Sigma"
	BAIKAL7_EXTENSION, BAIKAL8_EXTENSION, SIGMA_EXTENSION = "00", "xx", "bin"
	SIGMA_SECONDS_OFFSET = 2
	COMPONENTS_ORDER = "ZXY"
	MEMORY_BLOCK_SIZE = 120000
)

var BINARY_FILE_FORMATS = map[string]string{
	BAIKAL7_FMT: BAIKAL7_EXTENSION,
	BAIKAL8_FMT: BAIKAL8_EXTENSION,
	SIGMA_FMT: SIGMA_EXTENSION,
}


func truncate(num float64, precision uint8) float64 {
	return math.Round(num * math.Pow(10, float64(precision))) / math.Pow(10, float64(precision))
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


func getDatetimeStartSigma(dateNum uint32, timeNum uint32) (time.Time, error) {
	dateLine := strconv.FormatInt(int64(dateNum), 10)
	if len(dateLine) != 6 {
		return time.Time{}, BadHeaderData{"Invalid date in header"}
	}

	year, _ := strconv.ParseInt(dateLine[:2], 10, 64)
	year += 2000

	month, _ := strconv.ParseInt(dateLine[2:4], 10, 64)
	if month < 1 || month > 12 {
		return time.Time{}, BadHeaderData{"Invalid month in header"}
	}

	day, _ := strconv.ParseInt(dateLine[4:], 10, 64)

	firstMonthDay := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastMonthDay := firstMonthDay.AddDate(0, 1, 0).Add(-time.Nanosecond)
	if day < 1 || int(day) > lastMonthDay.Day() {
		return time.Time{}, BadHeaderData{"Invalid day in header"}
	}

	timeLine := fmt.Sprintf("%06d", int64(timeNum))
	hours, _ := strconv.ParseInt(timeLine[:2], 10, 64)
	if hours > 23 {
		return time.Time{}, BadHeaderData{"Invalid hours in header"}
	}

	minutes, _ := strconv.ParseInt(timeLine[2:4], 10, 64)
	if minutes > 59 {
		return time.Time{}, BadHeaderData{"Invalid minutes in header"}
	}

	seconds, _ := strconv.ParseInt(timeLine[4:], 10, 64)
	if seconds > 59 {
		return time.Time{}, BadHeaderData{"Invalid seconds in header"}
	}

	datetimeStart := time.Date(
		int(year), time.Month(month), int(day), 
		int(hours), int(minutes), int(seconds), 0, time.UTC)
	return datetimeStart, nil
}


func getCoordinatesSigma(longitudeLine string, latitudeLine string) (Coordinate, error) {
	if len(longitudeLine) != 9 {
		return Coordinate{}, BadHeaderData{"Invalid longitude in header"}
	}

	longitudeSymbol := longitudeLine[len(longitudeLine) - 1]
	if longitudeSymbol != 'E' && longitudeSymbol != 'W' {
		return Coordinate{}, BadHeaderData{"Invalid longitude in header"}
	}

	latitudeSymbol := latitudeLine[len(latitudeLine) - 1]
	if latitudeSymbol != 'N' && latitudeSymbol != 'S' {
		return Coordinate{}, BadHeaderData{"Invalid latitude in header"}
	}

	integerPart, _ := strconv.ParseFloat(longitudeLine[:3], 64)
	decimalPart, _ := strconv.ParseFloat(longitudeLine[3:len(longitudeLine) - 1], 64)
	longitude := truncate(integerPart + decimalPart / 60, 5)
	if longitudeSymbol == 'W' {
		longitude = -longitude
	}

	integerPart, _ = strconv.ParseFloat(latitudeLine[:2], 64)
	decimalPart, _ = strconv.ParseFloat(latitudeLine[2:len(latitudeLine) - 1], 64)
	latitude := truncate(integerPart + decimalPart / 60, 5)
	if latitudeSymbol == 'S' {
		latitude = -latitude
	}
	return Coordinate{longitude, latitude}, nil
}


func ReadBaikal7Header(path string) (FileHeader, error) {
	file, _ := os.Open(path)
	defer file.Close()

	channelsCount := UnsignedShortType{file, 0, 1}.convertToNumber()
	if int(channelsCount) != len(COMPONENTS_ORDER) {
		return FileHeader{}, BadHeaderData{message: "Invalid channels count"}
	}

	frequency := UnsignedShortType{file, 22, 1}.convertToNumber()

	srcCoords := DoubleType{file, 72, 2}.convertToArray()
	latitude, longitude := truncate(srcCoords[0], 5), truncate(srcCoords[1], 5)
	timeBegin := LongType{file, 104, 1}.convertToNumber()
	datetimeStart := getDatetimeStartBaikal7(timeBegin)
	return FileHeader{
		frequency: frequency, 
		datetimeStart: datetimeStart, 
		coordinate: Coordinate{
			Longitude: longitude, 
			Latitude: latitude}}, nil
}


func ReadBaikal8Header(path string) (FileHeader, error) {
	file, _ := os.Open(path)
	defer file.Close()

	channelsCount := UnsignedShortType{file, 0, 1}.convertToNumber()
	if int(channelsCount) != len(COMPONENTS_ORDER) {
		return FileHeader{}, BadHeaderData{message: "Invalid channels count"}
	}

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
	latitude, longitude := truncate(srcCoords[0], 5), truncate(srcCoords[1], 5)
	return FileHeader{
		frequency: frequency, 
		datetimeStart: datetimeStart, 
		coordinate: Coordinate{
			Longitude: longitude, 
			Latitude: latitude}}, nil 
}


func ReadSigmaHeader(path string) (FileHeader, error) {
	file, _ := os.Open(path)
	defer file.Close()

	channelsCount := UnsignedShortType{file, 12, 1}.convertToNumber()
	if int(channelsCount) != len(COMPONENTS_ORDER) {
		return FileHeader{}, BadHeaderData{message: "Invalid channels count"}
	}

	frequency := UnsignedShortType{file, 24, 1}.convertToNumber()
	latitudeSrc, longitudeSrc := CharType{file, 40, 8}.convert(), CharType{file, 48, 9}.convert()
	coordinates, err := getCoordinatesSigma(longitudeSrc, latitudeSrc)
	if err != nil {
		return FileHeader{}, err
	}

	datetimeSrc := UnsignedIntType{file, 60, 2}.convertToArray()
	datetimeStart, err := getDatetimeStartSigma(datetimeSrc[0], datetimeSrc[1])
	if err != nil {
		return FileHeader{}, err
	}
	
	return FileHeader{
		frequency: frequency, 
		datetimeStart: datetimeStart, 
		coordinate: coordinates} , nil
}


type BinaryFile struct {
	Path string
	ResampleFrequency uint16
	IsUseAvgValues bool
}

func (binFile BinaryFile) FileExtension() (string, error) {
	if len(binFile.Path) == 0 {
		return "", BadFilePath{message: "Empty file path"}
	}
	splitPath := strings.Split(binFile.Path, ".")
	return splitPath[len(splitPath) - 1], nil
}

func (binFile BinaryFile) FormatType() (string, error) {
	currentExtension, err := binFile.FileExtension()
	if err != nil {
		return "", err
	}

	for fileFormat, fileExtension := range BINARY_FILE_FORMATS {
		if currentExtension == fileExtension {
			return fileFormat, nil
		}
	}
	return "", BadFilePath{message: "Invalid file format"}
}

func (binFile BinaryFile) fileHeader() (FileHeader, error) {
	defaultValue := FileHeader{}
	formatType, err := binFile.FormatType()
	if err != nil {
		return defaultValue, err
	}

	switch formatType {
	case BAIKAL7_FMT:
		header, err := ReadBaikal7Header(binFile.Path)
		if err != nil {
			return defaultValue, err
		}
		return header, nil
	case BAIKAL8_FMT:
		header, err := ReadBaikal8Header(binFile.Path)
		if err != nil {
			return defaultValue, err
		}
		return header, nil
	case SIGMA_FMT:
		header, err := ReadSigmaHeader(binFile.Path)
		if err != nil {
			return defaultValue, err
		}
		return header, nil
	default:
		return defaultValue, BadFilePath{message: "Unknown format type"}
	}
}

func (binFile BinaryFile) GetResampleFrequency() (uint16, error) {
	header, err := binFile.fileHeader()
	if err != nil {
		return 0, err
	}

	switch freq := binFile.ResampleFrequency; {
	case freq < 0:
		return 0, InvalidResampleFrequency{message: fmt.Sprint(freq)}
	case freq == 0:
		return header.frequency, nil
	case header.frequency % freq == 0:
		return freq, nil
	default:
		return 0, InvalidResampleFrequency{message: fmt.Sprint(freq)}
	}
}

func (binFile BinaryFile) headerMemorySize() uint16 {
	return uint16(120 + 72 * len(COMPONENTS_ORDER))
}

func (binFile BinaryFile) discreteCount() (uint64, error) {
	info, err := os.Stat(binFile.Path)
	if err != nil {
		return 0, err
	}

	headerSize := binFile.headerMemorySize()

	size := info.Size()
	discreteCount := (uint64(size) - uint64(headerSize)) / (uint64(len(COMPONENTS_ORDER)) * uint64(UnsignedIntType{}.ByteSize()))
	return discreteCount, nil
}

func (binFile BinaryFile) secondsDuration() (float64, error) {
	header, err := binFile.fileHeader()
	if err != nil {
		return 0, err
	}

	discreteCount, err := binFile.discreteCount()
	if err != nil {
		return 0, err
	}

	frequency := header.frequency
	accuracy := uint8(math.Log10(float64(frequency)))

	deltaSeconds := truncate(float64(discreteCount) / float64(frequency), accuracy)
	return deltaSeconds, nil
}

func (binFile BinaryFile) DatetimeStart() (time.Time, error) {
	formatType, err := binFile.FormatType()
	if err != nil {
		return time.Time{}, err
	}

	header, err := binFile.fileHeader()
	if err != nil {
		return time.Time{}, err
	}

	offset := 0
	if formatType == SIGMA_FMT {
		offset = SIGMA_SECONDS_OFFSET
	}
	return header.datetimeStart.Add(time.Second * time.Duration(offset)), nil
}

func (binFile BinaryFile) DatetimeStop() (time.Time, error) {
	datetimeStart, err := binFile.DatetimeStart()
	if err != nil {
		return time.Time{}, err
	}

	secondsDuration, err := binFile.secondsDuration()
	if err != nil {
		return time.Time{}, err
	}

	nanosecondsDuration := secondsDuration * 1e9
	return datetimeStart.Add(time.Nanosecond * time.Duration(nanosecondsDuration)), nil
}

func (binFile BinaryFile) FileInfo() (FileInfo, error) {
	defaultValue := FileInfo{}

	formatType, err := binFile.FormatType()
	if err != nil {
		return defaultValue, err
	}

	path := binFile.Path

	header, err := binFile.fileHeader()
	if err != nil {
		return defaultValue, err
	}
	timeStart, _ := binFile.DatetimeStart()
	timeStop, err := binFile.DatetimeStop()
	if err != nil {
		return defaultValue, err
	}
	
	return FileInfo{
		Path: path, 
		FormatType: formatType, 
		Frequency: header.frequency, 
		TimeStart: timeStart, 
		TimeStop: timeStop, Coordinate: header.coordinate}, nil
	
}

func (binFile BinaryFile) IsGoodReadDatetimeStart(datetime time.Time) (bool, error) {
	datetimeStart, err := binFile.DatetimeStart()
	if err != nil {
		return false, err
	}

	datetimeStop, err := binFile.DatetimeStop()
	if err != nil {
		return false, err
	}

	secondsDiff := datetime.Sub(datetimeStart).Seconds()
	if secondsDiff < 0 {
		return false, InvalidDatetimeValue{message: "Reading time start is less than recording time start"}
	}

	secondsDiff = datetimeStop.Sub(datetime).Seconds()
	if secondsDiff <= 0 {
		return false, InvalidDatetimeValue{message: "Reading time start is more than recording time stop"}
	}

	return true, nil
}

func (binFile BinaryFile) IsGoodReadDatetimeStop(datetime time.Time) (bool, error) {
	datetimeStart, err := binFile.DatetimeStart()
	if err != nil {
		return false, err
	}

	datetimeStop, err := binFile.DatetimeStop()
	if err != nil {
		return false, err
	}

	secondsDiff := datetime.Sub(datetimeStart).Seconds()
	if secondsDiff <= 0 {
		return false, InvalidDatetimeValue{message: "Reading time stop is less than recording time start"}
	}

	secondsDiff = datetimeStop.Sub(datetime).Seconds()
	if secondsDiff < 0 {
		return false, InvalidDatetimeValue{message: "Reading time stop is more than recording time stop"}
	}

	return true, nil
}

func (binFile BinaryFile) resampleParameter() (uint16, error) {
	resampleFrequency, err := binFile.GetResampleFrequency()
	if err != nil {
		return 0, err
	}

	header, _  := binFile.fileHeader()
	return header.frequency / resampleFrequency, nil
}

func (binFile BinaryFile) componentIndex(component rune) (uint8, error) {
	for i, sourceComponent := range COMPONENTS_ORDER {
		if sourceComponent == component {
			return uint8(i), nil
		}
	}
	return 0, UnknownComponentName{message: string(component)}
}

func (binFile BinaryFile) getIndexesInterval(datetimeStart time.Time, datetimeStop time.Time) ([2]uint64, error) {
	defaultValue := [2]uint64{0, 0}
	_, err := binFile.IsGoodReadDatetimeStart(datetimeStart)
	if err != nil {
		return defaultValue, err
	}

	_, err = binFile.IsGoodReadDatetimeStop(datetimeStop)
	if err != nil {
		return defaultValue, err
	}

	resampleParameter, err := binFile.resampleParameter()
	if err != nil {
		return defaultValue, err
	}

	header, _ := binFile.fileHeader()
	originFrequency := header.frequency

	recordingDatetimeStart, _ := binFile.DatetimeStart()
	
	secondsDiff := datetimeStart.Sub(recordingDatetimeStart).Seconds()
	startIndex := uint64(math.Round(secondsDiff * float64(originFrequency)))

	secondsDiff = datetimeStop.Sub(recordingDatetimeStart).Seconds()
	stopIndex := uint64(math.Round(secondsDiff * float64(originFrequency)))

	signalLength := (stopIndex - startIndex) / uint64(resampleParameter)
	stopIndex = startIndex + signalLength * uint64(resampleParameter)

	return [2]uint64{startIndex, stopIndex}, nil
}
