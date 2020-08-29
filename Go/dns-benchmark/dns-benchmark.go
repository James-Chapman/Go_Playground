package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

var (
	totaltime   int64
	lookuptimes = make(map[string][]int64)
	rwm         sync.RWMutex
)

const (
	QUERIES = 500
)

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func get(key string) []int64 {
	rwm.RLock()
	defer rwm.RUnlock()
	return lookuptimes[key]
}

func appendx(server string, time int64) {
	rwm.Lock()
	defer rwm.Unlock()
	lookuptimes[server] = append(lookuptimes[server], time)
}

func lookupname(target string, server string, wg *sync.WaitGroup) {
	defer wg.Done()
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(target+".", dns.TypeA)
	r, t, err := c.Exchange(&m, server+":53")
	if err != nil {
		log.Println(err)
		return
	}

	totaltime += t.Nanoseconds()
	appendx(server, int64(t))
	//log.Printf("Took %v", t)
	if len(r.Answer) == 0 {
		log.Println("No results for " + target)
	}
}

func doesResultFileExist(name string) (bool, error) {
	matches, err := filepath.Glob(name)
	if err != nil {
		return false, err
	}
	return len(matches) > 0, nil
}

func main() {

	totaltime = 0

	txtfile, err := os.Open("top-1m.csv")
	checkError(err)
	defer txtfile.Close()
	buf, err := ioutil.ReadAll(txtfile)
	checkError(err)
	strbuf := string(buf)

	dnsfile, err := os.Open("dns.txt")
	checkError(err)
	defer dnsfile.Close()
	dnsbuf, err := ioutil.ReadAll(dnsfile)
	checkError(err)
	strdns := string(dnsbuf)
	servers := strings.Split(strdns, "\n")

	var lines = strings.Split(strbuf, "\n")
	line := " "
	domain := " "

	var wg sync.WaitGroup
	for s := 0; s < len(servers); s++ {
		log.Printf("Looking up against %s", servers[s])
		lookuptimes[servers[s]] = []int64{}
		for i := 0; i < QUERIES; i++ {
			line = lines[i]
			domain = strings.Split(line, ",")[1]
			wg.Add(1)
			go lookupname(domain, servers[s], &wg)
		}
		fmt.Printf("%d nanoseconds\n\n", totaltime)
		totaltime = 0
	}

	wg.Wait()

	t := time.Now()
	resultLine := ""

	for s := 0; s < len(servers); s++ {
		times := get(servers[s])
		var total int64
		var biggest int64
		var smallest int64
		total = 0
		biggest = times[0]
		smallest = times[0]

		for i := 0; i < len(times); i++ {
			total += times[i]
			if times[i] > biggest {
				biggest = times[i]
			}
			if times[i] < smallest {
				smallest = times[i]
			}
		}
		mean := total / int64(len(times))
		median := times[len(times)/2]
		range0 := biggest - smallest
		fmt.Printf("Mean (average) time taken for %s is .... %d nanoseconds\n", servers[s], mean)
		fmt.Printf("Median (middle) time taken for %s is ... %d nanoseconds\n", servers[s], median)
		fmt.Printf("Range for %s is ........................ %d nanoseconds\n\n", servers[s], range0)

		resultLine += fmt.Sprintf("%s,", t.Format("2006-01-02 15:04:05"))
		resultLine += fmt.Sprintf("%d,", mean)
		resultLine += fmt.Sprintf("%d,", median)
		resultLine += fmt.Sprintf("%d\n", range0)

		filename := fmt.Sprintf("Results_%s.csv", servers[s])

		exists, err := doesResultFileExist(filename)
		if err != nil {
			log.Println(err)
		}

		if !exists {
			newFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0660)
			if err != nil {
				log.Printf("Error opening %s\n", filename)
				continue
			}
			defer newFile.Close()
			newFile.WriteString("Date Time,Mean,Median,Range\n")
			newFile.Close()
		}

		resultFile, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0660)
		if err != nil {
			log.Printf("Error opening %s\n", filename)
			continue
		}
		defer resultFile.Close()

		resultFile.WriteString(resultLine)
	}
}
