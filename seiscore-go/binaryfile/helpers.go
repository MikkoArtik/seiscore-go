package binaryfile


import (
	"bytes"
	"encoding/binary"
	"os"
)


func readBinary(file *os.File, bytesCount uint16, skippingBytes uint16) []byte {
	bytes := make([]byte, bytesCount)
	file.ReadAt(bytes, int64(skippingBytes))
	return bytes
}


type CharType struct {
	filePointer *os.File
	skippingBytes uint16
	elementsCount uint16
}

func (dataType CharType) ByteSize() uint8 {
	return 1
}

func (dataType CharType) convert() string {
	bytesVal := readBinary(
		dataType.filePointer, uint16(dataType.ByteSize()) * dataType.elementsCount, 
		dataType.skippingBytes)
	return string(bytesVal)
}


type UnsignedShortType struct {
	filePointer *os.File
	skippingBytes uint16
	elementsCount uint16	
}

func (dataType UnsignedShortType) ByteSize() uint8 {
	return 2
}

func (dataType UnsignedShortType) getBytes() []byte {
	return readBinary(
		dataType.filePointer, uint16(dataType.ByteSize()) * dataType.elementsCount, 
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
	skippingBytes uint16
	elementsCount uint16
}

func (dataType UnsignedIntType) ByteSize() uint8 {
	return 4
}

func (dataType UnsignedIntType) getBytes() []byte {
	return readBinary(
		dataType.filePointer, uint16(dataType.ByteSize()) * dataType.elementsCount, 
		dataType.skippingBytes)
}

func (dataType UnsignedIntType) convertToNumber() uint32 {
	buffer := bytes.NewBuffer(dataType.getBytes())
	var result uint32
	binary.Read(buffer, binary.LittleEndian, &result)
	return result
}

func (dataType UnsignedIntType) convertToArray() []uint32 {
	buffer := bytes.NewBuffer(dataType.getBytes())
	result := make([]uint32, dataType.elementsCount)
	binary.Read(buffer, binary.LittleEndian, &result)
	return result
}


type DoubleType struct {
	filePointer *os.File
	skippingBytes uint16
	elementsCount uint16
}

func (dataType DoubleType) ByteSize() uint8 {
	return 8
}

func (dataType DoubleType) getBytes() []byte {
	return readBinary(
		dataType.filePointer, uint16(dataType.ByteSize()) * dataType.elementsCount, 
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
	skippingBytes uint16
	elementsCount uint16
}

func (dataType LongType) ByteSize() uint8 {
	return 8
}

func (dataType LongType) getBytes() []byte {
	return readBinary(
		dataType.filePointer, uint16(dataType.ByteSize()) * dataType.elementsCount, 
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
