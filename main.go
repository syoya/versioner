package main

import (
	"bufio"
	"fmt"
	version "github.com/hashicorp/go-version"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const filename string = "release_version"

var (
	show   = kingpin.Command("show", "Show the version.")
	create = kingpin.Command("create", "Create versioning file.")
	bump   = kingpin.Command("bump", "Bump version")
)

func main() {
	switch kingpin.Parse() {
	case create.FullCommand():
		if isVersionFile() {
			log.Fatalln("Cat't create versioning file. Versioning file already exists.")
		}
		content := []byte("v0.1.0")
		if err := ioutil.WriteFile(filename, content, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Versioning file was created.")
		return
	case show.FullCommand():
		content, err := readVersion()
		if err != nil {
			log.Fatalln("Can't open versioning file.", err)
		}
		fmt.Print(string(content))
		return
	case bump.FullCommand():
		content, err := readVersion()
		if err != nil {
			log.Fatalln("Can't open versioning file.", err)
		}
		v, err := version.NewVersion(string(content[1:]))
		if err != nil {
			log.Fatalln("Versioning format is invalid.", err)
		}
		fmt.Println("Version is", v.String())

		cv := v.Segments()
		cv[len(cv)-1] += 1
		nv := make([]string, len(cv))
		for i, n := range cv {
			nv[i] = strconv.Itoa(n)
		}
		newVersion := []byte("v" + strings.Join(nv, "."))
		fmt.Println("New version is", string(newVersion))

		os.Remove(filename)
		if err := ioutil.WriteFile(filename, newVersion, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Versioning file was updated.")
		return
	}
}

func readVersion() ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("There is no versioning file in this diretory. Please create versioning file.", err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	r := strings.NewReader(scanner.Text())
	return ioutil.ReadAll(r)
}

func isVersionFile() bool {
	_, err := os.Open(filename)
	if err != nil {
		return false
	}
	return true
}
