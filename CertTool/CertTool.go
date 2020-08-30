package main

import (
	"crypto/x509"
	"debug/pe"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	"go.mozilla.org/pkcs7"
)

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func toHexInt(n *big.Int) string {
	return fmt.Sprintf("%x", n) // or %X or upper case
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
		return nil, errors.New("Not signed PE file.")
	}

	f, err := os.Open(filePath)
	checkError(err)
	defer f.Close()

	buf = make([]byte, int64(size))
	f.ReadAt(buf, int64(vAddr+8))

	return buf, nil
}

// PrintCertDetails prints details of x509 Certificate
func PrintCertDetails(cert x509.Certificate) {
	fmt.Printf("Subject:        %s\n", cert.Subject)
	fmt.Printf("Issuer:         %s\n", cert.Issuer)
	fmt.Printf("Subject Key ID: %s\n\n", hex.EncodeToString(cert.SubjectKeyId))
}

func main() {
	var (
		inFilePtr = flag.String("in", "", "[Required] Path of PE file to open.")
		certFile  = flag.Bool("out", false, "[Optional] Write out PKCS7 cert to seperate file.")
	)

	flag.Parse()

	argLength := len(os.Args[1:])
	if argLength < 1 {
		flag.Usage()
		os.Exit(-1)
	}

	fi, err := os.Stat(*inFilePtr)
	buf, err := ExtractDigitalSignature(*inFilePtr)
	checkError(err)

	if *certFile == true {
		outFileName := fi.Name() + ".p7b"
		err = ioutil.WriteFile(outFileName, buf, 0644)
		checkError(err)
	}

	var pkcs7certPtr *pkcs7.PKCS7
	pkcs7certPtr, err = pkcs7.Parse(buf)
	checkError(err)

	//fmt.Print(pkcs7certPtr.Content)
	//fmt.Print(pkcs7certPtr.CRLs)
	//fmt.Print(pkcs7certPtr.Signers)
	//fmt.Print(pkcs7certPtr.Certificates)
	for _, pcert := range pkcs7certPtr.Certificates {
		PrintCertDetails(*pcert)
	}

}
