package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

func main() {
	fmt.Println("Listening dicom on port 11112...")
	storescp := exec.Command("./third_party/dcmtk/storescp", "-od", "./dicom", "-q", "-xcr", "printArg #f", "11112")
	// If main crash storescp still run and prevent main to restart
	storescpOut, err := storescp.StdoutPipe()
	if err != nil {
		log.Print(err)
	}
	if err := storescp.Start(); err != nil {
		log.Print(err)
	}

	scanner := bufio.NewScanner(storescpOut)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		s := fmt.Sprintf("./dicom/%s", m)
		dcm2jsonOut, err := exec.Command("./third_party/dcmtk/dcm2json", "-q", "-fc", s).Output()
		if err != nil {
			log.Print(err)
		}
		fmt.Println(string(dcm2jsonOut))

		if err := os.Remove(s); err != nil {
			log.Print(err)
		}
		dlp := gjson.GetBytes(dcm2jsonOut, "0040A730.Value.15.0040A730.Value.5.0040A730.Value.2.0040A300.Value.0.0040A30A.Value.0")
		ctdi := gjson.GetBytes(dcm2jsonOut, "0040A730.Value.15.0040A730.Value.5.0040A730.Value.0.0040A300.Value.0.0040A30A.Value")
		accessionNumber := gjson.GetBytes(dcm2jsonOut, "00080050.Value.0")

		id, err := uuid.NewUUID()
		if err != nil {
			log.Print(err)
		}
		message := []byte(fmt.Sprintf("{accessionNumber:%s,dlp:%s,ctdi:%s}", accessionNumber, dlp, ctdi))
		path := fmt.Sprintf("./json/%s.json", id)

		if err := ioutil.WriteFile(path, message, 0644); err != nil {
			log.Print(err)
		}
	}
	storescp.Wait()
}
