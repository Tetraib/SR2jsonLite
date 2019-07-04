package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	storescp := exec.Command("./third_party/dcmtk/storescp", "-od", "./dicom", "-q", "-xcr", "printArg #f", "11112")
	// dcm2jsonOut, err := exec.Command("date").Output()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	storescpOut, err := storescp.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := storescp.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(storescpOut)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)

	}
	storescp.Wait()
}
