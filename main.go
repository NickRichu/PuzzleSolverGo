// SiftTest project main.go
package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	//	"sort"
	"image"
	"image/draw"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
	//	"golang.org/x/tour/pic"
)

func main() {

	imageNames := readImagesIn("/home/nick/Documents/PetProjects/PuzzleSolverGo/pieces_aligned")
	guideImage := gocv.IMRead("hawaii_full.jpg", gocv.IMReadAnyColor)
	imageList := storeImagesInArray(imageNames)
	puzzleSolverWithoutRotation(guideImage, imageList)

}

func readImagesIn(path string) []string {
	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".jpg" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}

	return files
}

func storeImagesInArray(imageNames []string) []gocv.Mat {
	var imageList []gocv.Mat

	for _, file := range imageNames {
		imageList = append(imageList, gocv.IMRead(file, gocv.IMReadAnyColor))
	}

	return imageList
}

func showImage(windowName string, image gocv.Mat) {
	window := gocv.NewWindow(windowName)
	window.IMShow(image)
	gocv.WaitKey(0)
}

func puzzleSolverWithoutRotation(guideImage gocv.Mat, imageList []gocv.Mat) {
	var badImages []gocv.Mat
	//	var goodMatches []gocv.DMatch
	imageGuide, _ := guideImage.ToImage()
	mockCanvas := image.Rectangle{image.Point{0, 0}, imageGuide.Bounds().Max}
	trueCanvas := image.NewRGBA(mockCanvas)

	var finalCanvas gocv.Mat
	mask := gocv.NewMat()
	si := contrib.NewSIFT()
	bf := gocv.NewBFMatcher()
	kpGuide, desc1 := si.DetectAndCompute(guideImage, mask)

	for _, singleImage := range imageList {
		_, desc2 := si.DetectAndCompute(singleImage, mask)
		if !desc2.Empty() {
			match := bf.KnnMatch(desc1, desc2, 2)
			goodMatches := goodMatchFilter(match)

			for _, match := range goodMatches[:1] {
				currentKP1 := kpGuide[match.QueryIdx]
				locationx := int(math.RoundToEven(currentKP1.X)/50) * 50
				locationy := int(math.RoundToEven(currentKP1.Y)/50) * 50

				imageSingleImage, _ := singleImage.ToImage()
				locationPoints := image.Point{X: locationx, Y: locationy}

				currentImage := image.Rectangle{locationPoints, locationPoints.Add(imageSingleImage.Bounds().Size())}
				fmt.Println(currentImage.Max)
				draw.Draw(trueCanvas, currentImage, imageSingleImage, image.Point{0, 0}, draw.Src)
				finalCanvas, _ = ToRGB8(trueCanvas)
				print("----------------------------------\n")
				print(locationx, " coordinate x \n")
				print(locationy, " coordinate y\n\n")
			}

		} else {
			badImages = append(badImages, singleImage)

			print("it hit this")
			break
		}

	}
	showImage("canvas", finalCanvas)

}

func goodMatchFilter(match [][]gocv.DMatch) []gocv.DMatch {
	var goodMatches []gocv.DMatch
	for i := 0; i < len(match); i++ {
		if match[i][0].Distance < 0.5*match[i][1].Distance {
			//			fmt.Print("zero value:", match[i][0])
			goodMatches = append(goodMatches, match[i][0])
			return goodMatches
		}
	}
	return goodMatches
}

func ToRGB8(img image.Image) (gocv.Mat, error) {
	bounds := img.Bounds()
	x := bounds.Dx()
	y := bounds.Dy()
	bytes := make([]byte, 0, x*y*3)

	//don't get surprised of reversed order everywhere below
	for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
		for i := bounds.Min.X; i < bounds.Max.X; i++ {
			r, g, b, _ := img.At(i, j).RGBA()
			bytes = append(bytes, byte(b>>8), byte(g>>8), byte(r>>8))
		}
	}
	return gocv.NewMatFromBytes(y, x, gocv.MatTypeCV8UC3, bytes)
}
