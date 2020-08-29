/**
*
* Search locations in dirs.txt recursively for files in files.txt and copy to designated location
*
**/

package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	debugLog bool
	dstDir   string
)

func checkError(e error, fatal bool) {
	if e != nil {
		if fatal == true {
			log.Fatal(e)
		} else {
			log.Printf("[WARNING] %s", e)
		}
	}
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func readFile(filePath string) string {
	hFile, err := os.Open(filePath)
	checkError(err, true)
	defer hFile.Close()
	buf, err := ioutil.ReadAll(hFile)
	checkError(err, true)
	return string(buf)
}

func writeFile(filePath string, data []byte, append bool) {
	var hFile *os.File
	// if file does not exist, create it
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		hFile, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0660)
		checkError(err, true)
	} else {
		if append {
			hFile, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0660)
			checkError(err, true)
		} else {
			hFile, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0660)
			checkError(err, true)
		}
	}

	bytesWritten, err := hFile.Write(data)
	checkError(err, true)
	err = hFile.Close()
	checkError(err, true)
	if debugLog {
		log.Printf("%d bytes written to %s", bytesWritten, filePath)
	}
}

func generateMD5(filePath string) (string, error) {
	var md5String string
	file, err := os.Open(filePath)
	if err != nil {
		return md5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return md5String, err
	}
	hashBytes := hash.Sum(nil)[:16]
	md5String = hex.EncodeToString(hashBytes)
	return md5String, nil
}

func copyFile(src string, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func processFiles(dirData string, fileData string) {
	dirData = strings.ReplaceAll(dirData, "\r", "")
	fileData = strings.ReplaceAll(fileData, "\r", "")

	dirLines := strings.Split(dirData, "\n")
	fileLines := strings.Split(fileData, "\n")

	if debugLog {
		log.Println(dirLines)
		log.Println(fileLines)
	}

	var fileMap = make(map[string]int)

	for _, line := range fileLines {
		if line == "" {
			continue
		}

		if debugLog {
			log.Printf("line: %s", line)
		}

		fileMap[line] = 1
	}

	for _, line := range dirLines {
		if line == "" {
			continue
		}

		if debugLog {
			log.Printf("line: %s", line)
		}

		err := filepath.Walk(line, func(pathx string, info os.FileInfo, err error) error {

			// first thing to do, check error. and decide what to do about it
			if err != nil {
				log.Printf("Error %v\n", err)
			} else {
				// If not a dir, i.e. if a file.
				if info.IsDir() == false {
					fmt.Print(".")
					if _, ok := fileMap[info.Name()]; ok {
						if debugLog {
							fmt.Printf("  dir: %v\n", filepath.Dir(pathx))
							fmt.Printf("  file name %v\n", info.Name())
							fmt.Printf("  extenion: %v\n", filepath.Ext(pathx))
						}

						var srcFilePath = path.Join(filepath.Dir(pathx), info.Name())
						if debugLog {
							log.Printf("srcFilePath: %s", srcFilePath)
						}

						md5String, _ := generateMD5(srcFilePath)
						if debugLog {
							log.Printf("MD5 for %s: %s", srcFilePath, md5String)
						}

						var dstFilePath = path.Join(dstDir, md5String, info.Name())
						if debugLog {
							log.Printf("dstFilePath: %s", dstFilePath)
						}
						if _, err := os.Stat(path.Join(dstDir, md5String)); os.IsNotExist(err) {
							os.MkdirAll(path.Join(dstDir, md5String), os.ModePerm)
						}
						copyFile(srcFilePath, dstFilePath)
					}
				}
			}

			return nil
		})
		checkError(err, false)
	}
}

func main() {
	var (
		dirFile   = flag.String("dirs", "dirs.txt", "File with dirs to search.")
		fileFile  = flag.String("files", "files.txt", "File with files to save.")
		outputDir = flag.String("out", "copied", "Dir in which to place copied files.")
		debugFlag = flag.Bool("debug", false, "Enable debug logging")
	)

	flag.Parse()
	debugLog = *debugFlag
	dstDir = *outputDir

	SetupCloseHandler()

	dirData := readFile(*dirFile)
	fileData := readFile(*fileFile)

	processFiles(dirData, fileData)

}
