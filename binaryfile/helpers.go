package binaryfile


import (
	"bytes"
	"encoding/binary"
	"os"
)


func readBinary(file *os.File, bytesCount int, skippingBytes int) []byte {
	bytes := make([]byte, bytesCount)
	file.ReadAt(bytes, int64(skippingBytes))
	return bytes
}


type CharType struct {
	filePointer *os.File
	skippingBytes int
	elementsCount int
}

func (dataType CharType) byteSize() int {
	return 1
}

func (dataType CharType) convert() string {
	bytesVal := readBinary(
		dataType.filePointer, dataType.byteSize() * dataType.elementsCount, 
		dataType.skippingBytes)
	return string(bytesVal)
}


type UnsignedShortType struct {
	filePointer *os.File
	skippingBytes int
	elementsCount int	
}

func (dataType UnsignedShortType) byteSize() int {
	return 2
}

func (dataType UnsignedShortType) getBytes() []byte {
	return readBinary(
		dataType.filePointer, dataType.byteSize() * dataType.elementsCount, 
		dataType.skippingBytes)
}

func (dataType UnsignedShortType) convertToNumber() uint16 {
	buffer := bytes.NewBuffer(dataType.getBytes())
	var result uint16
	binary.Read(buffer, binary.LittleEndian, &result)
	return result
}

func (dataType UnsignedShortType) convertToArray() []uint16 {
	buffer := bytes.NewBuffer(dataType.getBytes())
	result := make([]uint16, dataType.elementsCount)
	binary.Read(buffer, binary.LittleEndian, &result)
	return result
}


type UnsignedIntType struct {
	filePointer *os.File
	skippingBytes int
	elementsCount int
}

func (dataType UnsignedIntType) byteSize() int {
	return 4
}

func (dataType UnsignedIntType) getBytes() []byte {
	return readBinary(
		dataType.filePointer, dataType.byteSize() * dataType.elementsCount, 
		dataType.skippingBytes)
}

func (dataType UnsignedIntType) convertToNumber() int32 {
	buffer := bytes.NewBuffer(dataType.getBytes())
	var result int32
	binary.Read(buffer, binary.LittleEndian, &result)
	return result
}

func (dataType UnsignedIntType) convertToArray() []int32 {
	buffer := bytes.NewBuffer(dataType.getBytes())
	result := make([]int32, dataType.elementsCount)
	binary.Read(buffer, binary.LittleEndian, &result)
	return result
}


type DoubleType struct {
	filePointer *os.File
	skippingBytes int
	elementsCount int
}

func (dataType DoubleType) byteSize() int {
	return 8
}

func (dataType DoubleType) getBytes() []byte {
	return readBinary(
		dataType.filePointer, dataType.byteSize() * dataType.elementsCount, 
		dataType.skippingBytes)
}

func (dataType DoubleType) convertToNumber() float64 {
	buffer := bytes.NewBuffer(dataType.getBytes())
	var result float64
	binary.Read(buffer, binary.LittleEndian, &result)
	return result
}

func (dataType DoubleType) convertToArray() []float64 {
	buffer := bytes.NewBuffer(dataType.getBytes())
	result := make([]float64, dataType.elementsCount)
	binary.Read(buffer, binary.LittleEndian, &result)
	return result
}


type LongType struct {
	filePointer *os.File
	skippingBytes int
	elementsCount int
}

func (dataType LongType) byteSize() int {
	return 8
}

func (dataType LongType) getBytes() []byte {
	return readBinary(
		dataType.filePointer, dataType.byteSize() * dataType.elementsCount, 
		dataType.skippingBytes)
}

func (dataType LongType) convertToNumber() uint64 {
	buffer := bytes.NewBuffer(dataType.getBytes())
	var result uint64
	binary.Read(buffer, binary.LittleEndian, &result)
	return result
}

func (dataType LongType) convertToArray() []uint64 {
	buffer := bytes.NewBuffer(dataType.getBytes())
	result := make([]uint64, dataType.elementsCount)
	binary.Read(buffer, binary.LittleEndian, &result)
	return result
}
