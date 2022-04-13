package main


import (
	"log"
	"example.com/binaryfile"
)



func main() {
	a := binaryfile.ReadBaikal8Header("/media/michael/Data/Projects/GraviSeismicComparation/Vibrostend/seismic/20-01-2022/K07_2022-01-20_08-21-05/K07_2022-01-20_08-21-05.xx")
	log.Println(a)
	b := binaryfile.ReadBaikal7Header("/media/michael/Data/Projects/HydroFracturing/DemkinskoeDeposit/Demkinskoe_4771/Binary/HF_0019_2019-08-16_08-31-08_90041_527.00")
	log.Println(b)

	c := binaryfile.ReadSigmaHeader("/media/michael/Data/Projects/GraviSeismicComparation/Vibrostend/seismic/19-01-2022/SigmaN011_2022-01-19_10-06-11.bin")
    log.Println(c)
}