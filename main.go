package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

var yellowBg = color.New(color.FgBlack, color.BgYellow).PrintFunc()
var cyanBg = color.New(color.FgBlack, color.BgCyan).PrintFunc()

func main() {
	color.HiGreen("╔══════════════════════════════════════╗")
	color.HiGreen("║                                      ║")
	color.HiGreen("╟── MATIC NI SYNCAL FOTOAPARATOV APP ──╢")
	color.HiGreen("║                                      ║")
	color.HiGreen("╚══════════════════════════════════════╝")
	filePaths := getRelativePaths([]string{".jpg", ".png"})
	fmt.Printf("Program je nasel %d datotek (.jpg in .png)\n\n", len(filePaths))

	rawOffset, err := readLine(" >> Vnesi za koliko zelis zamakniti vsako sliko (npr. -1h):")
	checkErr(err)

	parsedOffset, err := time.ParseDuration(rawOffset)
	checkErr(err)

	fmt.Printf("Vnesel si %s\n", parsedOffset.String())
	confirmation, err := readLine(" >> POTRDI (y/N):")
	checkErr(err)

	fmt.Println()

	if strings.ToLower(confirmation) != "y" {
		return
	}

	for _, filePath := range filePaths {
		// open file
		file, err := os.Open(filePath)
		checkErr(err)
		defer file.Close()

		// get file info
		fileStat, err := file.Stat()
		checkErr(err)

		modTime := fileStat.ModTime()
		newTime := modTime.Add(parsedOffset)

		color.Cyan("----------------------------")
		yellowBg(" Slika: ")
		color.Yellow(" " + filePath)

		cyanBg("  Prej: ")
		color.Cyan(" " + modTime.String()[:19])

		cyanBg(" Potem: ")
		color.Cyan(" " + newTime.String()[:19])

		os.Chtimes(filePath, newTime, newTime)
	}

	color.Cyan("----------------------------")

	fmt.Println()
	cyanBg(" --- KONČANO --- ")
	fmt.Scanln() // pause
}

// read whole sentence from stdin
func readLine(pre string) (string, error) {
	color.New(color.FgHiGreen).Print(pre + " ")
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	return string(line), nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func getRelativePaths(allowedExtensions []string) []string {
	var filePaths []string

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		// check file extension
		if !contains(allowedExtensions, filepath.Ext(path)) {
			return nil
		}

		// append file to a file list
		filePaths = append(filePaths, path)
		return nil
	})
	checkErr(err)
	return filePaths
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
