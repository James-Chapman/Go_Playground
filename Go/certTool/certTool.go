package main

import (
	//"crypto/x509"
	"debug/pe"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	//"encoding/pem"
	"github.com/fullsailor/pkcs7"
)

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// ExtractDigitalSignature extracts a digital signature specified in a signed PE file.
// It returns a digital signature (pkcs#7) in bytes.
func ExtractDigitalSignature(filePath string) (buf []byte, err error) {
	pefile, err := pe.Open(filePath)
	checkError(err)
	defer pefile.Close()

	var vAddr uint32
	var size uint32
	switch t := pefile.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		vAddr = t.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_SECURITY].VirtualAddress
		size = t.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_SECURITY].Size
	case *pe.OptionalHeader64:
		vAddr = t.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_SECURITY].VirtualAddress
		size = t.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_SECURITY].Size
	}

	if vAddr <= 0 || size <= 0 {
		return nil, errors.New("Not signed PE file")
	}

	f, err := os.Open(filePath)
	checkError(err)
	defer f.Close()

	buf = make([]byte, int64(size))
	f.ReadAt(buf, int64(vAddr+8))

	return buf, nil
}

func main() {

	filePtr := flag.String("file", "", "Path of PE file to open.")
	flag.Parse()

	buf, err := ExtractDigitalSignature(*filePtr)
	checkError(err)

	err = ioutil.WriteFile("extracted_cert.p7b", buf, 0644)
	checkError(err)
	
	pkcs7certPtr, err := pkcs7.Parse(buf)
	checkError(err)

	fmt.Printf("%#v\n", *pkcs7certPtr)
	fmt.Print(pkcs7certPtr.Content)
	fmt.Print(pkcs7certPtr.Certificates)
	fmt.Print(pkcs7certPtr.CRLs)
	fmt.Print(pkcs7certPtr.Signers)
}
