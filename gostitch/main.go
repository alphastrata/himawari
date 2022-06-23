package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
	"gocv.io/x/gocv"
)

//=============================================
//        USER INPUT
//=============================================
type DateTime struct {
	year, month, day, hhmmss uint
}

// parses the args to stitch the appropriate tiles, from the approtriate date
func parseParamsFromSYSArgs() DateTime {
	args := os.Args
	fmt.Println(args)
	if len(args) != 3 {
		fmt.Println("Usage: gostitch <YYYYMMDD> <HHMMSS>")
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

	return DateTime{
		year:   uint(year),
		month:  uint(month),
		day:    uint(day),
		hhmmss: uint(hhmmss),
	}
}

//=============================================
//        IMAGES
//=============================================
// helper func to show an image on screen, to help with debugging
func showImage(img gocv.Mat) {
	window := gocv.NewWindow("himawari")

	window.IMShow(img)

	window.WaitKey(0)
	if window.WaitKey(1) >= 0 {
		window.Close()
		return
	}
}

func loadImage(path string) gocv.Mat {
	img := gocv.IMRead(path, gocv.IMReadColor)
	if img.Empty() {
		log.Printf("Error reading image from: %v\n", path)
		return gocv.Mat{}
	}
	return img
}

func saveImage(filename string, img gocv.Mat) bool {
	gocv.IMWrite(filename, img)
	log.Println("Processed: ", filename)
	return true
}

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

//=============================================
//        FILES
//=============================================
func readinTiles(dir string) []Tile {
	if !exists(dir) {
		log.Fatalf("Error: directory %v does not exist\n", dir)
		return nil
	}
	files, e := ioutil.ReadDir(dir)
	CheckError(e)

	var tiles []Tile

	for _, f := range files {
		newname := f.Name()

		c := strings.Replace(strings.Split(newname, "_")[1], ".png", "", -1)
		c = (strings.Replace(c, "C", "", -1))
		col, e := strconv.Atoi(c)
		CheckError(e)

		r := strings.Split(newname, "_")[0][6:]
		r = (strings.Replace(r, "R", "", -1))
		row, e := strconv.Atoi(r)
		CheckError(e)

		tiles = append(tiles,
			Tile{
				img:  loadImage(fmt.Sprintf("%s/%s", dir, f.Name())),
				path: fmt.Sprintf("%s/%s", dir, f.Name()),
				name: f.Name(),
				row:  row,
				col:  col,
			})

	}
	return tiles
}
func exists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

//=============================================
//        TILE
//=============================================
type Tile struct {
	img  gocv.Mat
	name string
	path string
	row  int
	col  int
}
type TilePos struct {
	row int
	col int
}

func TilesToMap(t []Tile) map[TilePos]gocv.Mat {
	// build a map of tiles, Key = row, col with Value = Tile
	m := make(map[TilePos]gocv.Mat)
	for idx := 0; idx < 20; idx++ {
		if t[idx].row == idx {
			for jdx := 0; jdx < 400; jdx++ {
				m[TilePos{t[jdx].row, t[jdx].col}] = t[jdx].img
			}
		}

	}

	return m
}

//=============================================
//        CONCATENATORS
//=============================================
// helper func to build vertical stacks, src1 is on the top
func vJoin(src1, src2 gocv.Mat) gocv.Mat {
	mat := gocv.NewMat()
	gocv.Vconcat(src1, src2, &mat)
	return mat
}

// helper func to build horizontal stacks, src1 goes on the lhs
func hJoin(src1, src2 gocv.Mat) gocv.Mat {
	mat := gocv.NewMat()
	gocv.Hconcat(src1, src2, &mat)
	return mat
}

// helper func to build horizontal stacks, from a map of tiles
func hConcat(m map[TilePos]gocv.Mat) []gocv.Mat {
	strips := []gocv.Mat{}
	img := gocv.NewMat()

	for r := 0; r < 20; r++ {
		img = m[TilePos{0, 0 + r}]
		for c := 1; c < 20; c++ {
			img = hJoin(img, m[TilePos{c, r}])

		}
		strips = append(strips, img)

	}
	return strips
}

// helper func to build vertical stacks, from a slices of tiles
func vConcat(stack []gocv.Mat) gocv.Mat {
	img := stack[0]

	for i := 0; i < len(stack)-1; i++ {
		img = vJoin(img, stack[i+1])
	}
	return img
}

//=============================================
//       Full Discs
//=============================================
// TODO: Make the progressbar work in here, pass it in, pass refs to it in v and h concat so they can all tick it.
func processFullDiscs(dir, filepath_out string, show bool) {
	tiles := readinTiles(dir)
	m := TilesToMap(tiles)

	// Stack tiles into a full disc
	hstack := hConcat(m)
	full_disc := vConcat(hstack)

	// Save full discs
	saveImage(filepath_out, full_disc)
	log.Println("Processed: ", filepath_out)
	hstack = nil

	// Show the full_disc if so inclined, useful for debugging
	if show {
		showImage(full_disc)
	}

}

//=============================================
//        STITCH THOSE TILES!!!
//=============================================
func main() {

	dt := parseParamsFromSYSArgs()

	log.Printf("On DAY: %02d Of MONTH: %d at HOUR: %d", dt.day, dt.month, dt.hhmmss)
	bar := progressbar.NewOptions(100,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("[cyan][1/x][reset] Stitching tiles "),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	// Note the location your terminal is at when you run this seems to be important for these relative ../../s to work.
	filename_dir := filepath.Join("../", "completed", fmt.Sprintf("%d%02d%02d/", dt.year, dt.month, dt.day))
	filename_out := filepath.Join(filename_dir, fmt.Sprintf("/fulldisc-%d-%02d-%d %06d.png", dt.year, dt.month, dt.day, dt.hhmmss/100))
	tile_dir := filepath.Join("../", "tiles", fmt.Sprintf("%d%02d%02d/%06d", dt.year, dt.month, dt.day, dt.hhmmss))

	// NOTE: Blocking and concurrent versions of this func are available
	processFullDiscs(tile_dir, filename_out, false)
	log.Println("Done!")
	bar.Finish()

}
