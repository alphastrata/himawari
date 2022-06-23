package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

//=============================================
//        USER INPUT
//=============================================
// parses the arguments passed to himarawi which become the date and time to fetch
func parseParamsFromSYSArgs() UserInput {
	args := os.Args
	fmt.Println(args)
	if len(args) != 3 {
		fmt.Println("Usage: himawari <YYYYMMDD> <HHMMSS>")
		fmt.Println("Note : the HHMMSS are values from 0 to 160000, where 0 is midday.\nAs the MM values represent minutes, they can never be more than 50.\nThe seconds values will ALWAYS be 00.")
		os.Exit(1)
	}
	y := args[1][:4]
	year, _ := strconv.Atoi(y)
	m := args[1][4:6]
	month, _ := strconv.Atoi(m)
	d := args[1][6:8]
	day, _ := strconv.Atoi(d)

	hms := args[2]
	hhmmss, _ := strconv.Atoi(hms)

	log.Println("Fetching for :", year, month, day, hhmmss)

	return UserInput{
		year:   uint(year),
		month:  uint(month),
		day:    uint(day),
		hhmmss: uint(hhmmss),
	}
}

// UserInput is a convinience struct for managing the dates and times, as they're everywhere
type UserInput struct {
	year, month, day, hhmmss uint
}

//=============================================
//        FILES
//=============================================
type Folders struct {
	tiles     string
	completed string
}

func exists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

// builds the filename for the saved tiles
func buildTileFilename(url, sixnum, ymd string) string {
	url_array := strings.Split(url, "/")

	filename := (url_array[len(url_array)-1])
	filename = strings.Replace(filename, "_", "R", 1)
	filename = strings.Replace(filename, "_", "_C", 1)

	return "tiles/" + ymd + "/" + sixnum + "/" + filename
}

// helper function to count the number of tiles successfully downloaded,
// useful for when you're resuming or going back over a range that may be
// incomplete.
func countTiles(folders Folders) int {
	files, _ := filepath.Glob(folders.tiles + "/*")
	return len(files)
}

// creates the directories stucture for the tiles and the complete disc images,
// important because keeping files organised is important
func setupDirs(file_ymd string, sixnum uint) Folders {
	tiles := fmt.Sprintf("tiles/%s/%06d", file_ymd, sixnum)
	completed := fmt.Sprintf("completed/%s/", file_ymd)
	if !exists(tiles) {
		os.MkdirAll(tiles, 0777)
	}

	if !exists(completed) {
		os.MkdirAll(completed, 0777)
	}
	return Folders{tiles, completed}
}

//=============================================
//        URLS
//=============================================
func downloadFromURL(url string, filename string, wg *sync.WaitGroup) {
	defer wg.Done() // ensures the corouties won't race

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer response.Body.Close()

	// Create output file
	outFile, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer outFile.Close()

	// Copy data from HTTP response to file
	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// builds the url for download_from_url()
func URLBuilder(base, d_value, ymd, sixnum string, row, col int) string {
	url := fmt.Sprintf("%s%sd/550/%s/%s_%d_%d.png", base, d_value, ymd, sixnum, col, row)
	return url
}

// builds the ymd string for the urls
func ymdBuilder(y, m, d uint) string {
	return fmt.Sprintf("%d/%02d/%02d", y, m, d)
}

//=============================================
//        TILES
//=============================================
func executeTiles(file_ymd, url_ymd, d_value, base string, sixnum string) {
	bar := progressbar.NewOptions(400,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("[cyan][%][reset] Fetching tiles "),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	wg := new(sync.WaitGroup)

	COLMAX := 19
	ROWMAX := 19
	for row := 0; row <= ROWMAX; row++ {
		for col := 0; col <= COLMAX; col++ {
			url := URLBuilder(base, d_value, url_ymd, sixnum, row, col)
			filename := buildTileFilename(url, sixnum, file_ymd)
			if !exists(filename) {
				wg.Add(1)
				go downloadFromURL(url, filename, wg)
			}
			bar.Add(1)
		}
	}
	wg.Wait()
	bar.Finish()
}

// fetch/download the tiles, from the built urls, for the date and time specified
func fetchTileHelper(year, month, day, sixnum uint) Folders {
	d_value := "20"
	url_ymd := ymdBuilder(year, month, day)
	file_ymd := strings.ReplaceAll(url_ymd, "/", "")
	base := "https://himawari8.nict.go.jp/img/D531106/"

	hhmmdd := fmt.Sprintf("%06d", sixnum)

	folders := setupDirs(file_ymd, sixnum)

	executeTiles(file_ymd, url_ymd, d_value, base, hhmmdd)
	return folders
}

// Helper function to start the scraping process.
func scrape(year, month, day, hhmmss uint) bool {

	folders := fetchTileHelper(year, month, day, hhmmss)
	//time.Sleep(time.Second)

	log.Printf("Tiles gathered : %d\n", countTiles(folders))
	log.Println("Tile located  : " + folders.tiles)

	return true
}

//=============================================
//        FULL-DISC
//=============================================
func buildFullDisc(sa UserInput, tool string) bool {
	var runwith string

	ymd := fmt.Sprintf("%04d%02d%02d", sa.year, sa.month, sa.day)
	hhmmss := fmt.Sprintf("%06d", sa.hhmmss)

	// if tool == "go" {
	// 	runwith = fmt.Sprintf("cd gostitch;go run stitchers/gostitch/main.go %s %s", ymd, hhmmss)
	// 	fmt.Println("RANWITH: ", runwith)
	// 	runStitcher(sa, ymd, hhmmss, runwith)
	// 	return true
	// }
	if tool == "python" {
		runwith = fmt.Sprintf("python scripts/stitcher.py %s %s", ymd, hhmmss)
		fmt.Println("RANWITH: ", runwith)
		runStitcher(runwith)
		return true
	} else {
		log.Print("Invalid tool, valid options: python, go")
		return false
	}

}
func runStitcher(runwith string) bool {
	cmd := exec.Command(runwith)
	out, err := cmd.Output()

	if err != nil {
		log.Println(err)
		log.Println(cmd)
		log.Println(string(out))
		return false
	}

	return true
}

//=============================================
//        MAIN
//=============================================
func main() {
	start := time.Now()
	fmt.Println("Himawari")
	fmt.Println(os.Getwd())

	sa := parseParamsFromSYSArgs()
	scrape(sa.year, sa.month, sa.day, sa.hhmmss)
	buildFullDisc(sa, "python") //NOTE: The go one has been removed from this version -- the gocv module isn't stable enough for this application...

	elapsed := time.Since(start)

	fmt.Printf("Himawari Runtime: %s\n", elapsed)
	fmt.Printf(".......................................................................................\n")

}
