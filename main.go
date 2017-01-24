package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/deckarep/golang-set"
	"gopkg.in/gographics/imagick.v2/imagick"
	//	"./color"
	"github.com/toduq/ciedecompress/color"
)

type ProcessingPixel struct {
	X        int
	Y        int
	R        float64
	G        float64
	B        float64
	Lab      color.Lab
	Similars mapset.Set
}

var (
	srcFile  string
	destFile string
	size     int
	border   float64
	workers  int
	verbose  bool
)

func main() {
	flag.StringVar(&srcFile, "i", "src.bmp", "src file")
	flag.StringVar(&destFile, "o", "dest.bmp", "dest file")
	flag.IntVar(&size, "size", 8, "squash rect size")
	flag.Float64Var(&border, "border", 1.0, "squash diff limit")
	flag.IntVar(&workers, "workers", 3, "worker size")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.Parse()

	if verbose {
		fmt.Println(srcFile, "to", destFile, "with {size:", size, ", border: ", border, "}")
	}

	imagick.Initialize()
	defer imagick.Terminate()
	img := imagick.NewMagickWand()
	_, err := os.Stat(srcFile)
	if err != nil {
		panic(err)
	}
	err = img.ReadImage(srcFile)
	if err != nil {
		panic(err)
	}

	width := int(img.GetImageWidth())
	height := int(img.GetImageHeight())

	cells := (width / size) * (height / size)
	processed := 0

	jobs := make(chan []ProcessingPixel)
	pipe := make(chan []ProcessingPixel)
	done := make(chan bool)

	diffWorker := func(jobs <-chan []ProcessingPixel, pipe chan<- []ProcessingPixel) {
		for pixels := range jobs {
			for from, pix1 := range pixels {
				for to, pix2 := range pixels {
					if from <= to {
						continue
					}
					diff := pix1.Lab.Diff(pix2.Lab)
					if diff > border {
						continue
					}
					pix1.Similars.Add(to)
					pix2.Similars.Add(from)
				}
			}
			pipe <- pixels
		}
	}
	replaceWorker := func(pipe <-chan []ProcessingPixel, done chan<- bool) {
		for pixels := range pipe {
			if pixels == nil {
				done <- true
				return
			}
			for true {
				maxIndex, maxValue := -1, 0
				for index, pix := range pixels {
					length := pix.Similars.Cardinality()
					if length > maxValue {
						maxIndex = index
						maxValue = length
					}
				}
				if maxIndex == -1 {
					break
				}
				copyFrom := maxIndex
				copyTos := pixels[maxIndex].Similars
				for copyToInterface := range copyTos.Iterator().C {
					copyTo, _ := copyToInterface.(int)
					rgb := []float64{pixels[copyFrom].R, pixels[copyFrom].G, pixels[copyFrom].B}
					img.ImportImagePixels(pixels[copyTo].X, pixels[copyTo].Y, 1, 1, "RGB", imagick.PIXEL_DOUBLE, rgb)
				}
				deleteList := copyTos.Clone()
				deleteList.Add(copyFrom)
				for index, pix := range pixels {
					pixels[index].Similars = pix.Similars.Difference(deleteList)
				}
			}
		}
	}

	for j := 0; j < workers; j++ {
		go diffWorker(jobs, pipe)
	}
	go replaceWorker(pipe, done)

	for i := 0; i < cells; i++ {
		x := (i % (width / size)) * size
		y := (i / (height / size)) * size
		if processed%100 == 0 && verbose {
			fmt.Println("processing", cells, "->", processed)
		}
		processed++
		pixels := make([]ProcessingPixel, 0, size*size)
		for dy := 0; dy < size; dy++ {
			for dx := 0; dx < size; dx++ {
				pix, _ := img.GetImagePixelColor(x+dx, y+dy)
				r, g, b := pix.GetRed(), pix.GetGreen(), pix.GetBlue()
				lab := ProcessingPixel{X: x + dx, Y: y + dy, R: r, G: g, B: b, Lab: color.FromRgb(r, g, b), Similars: mapset.NewSet()}
				pixels = append(pixels, lab)
			}
		}
		jobs <- pixels
	}
	close(jobs)
	pipe <- nil
	<-done
	img.WriteImage(destFile)
}
