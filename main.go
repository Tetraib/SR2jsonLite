package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
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
		fmt.Println(m)
		s := fmt.Sprintf("./dicom/%s", m)
		dcm2jsonOut, err := exec.Command("./third_party/dcmtk/dcm2json", s).Output()
		if err != nil {
			log.Print(err)
		}
		if err := os.Remove(s); err != nil {
			log.Print(err)
		}
		fmt.Println(string(dcm2jsonOut))

	}
	storescp.Wait()
}
