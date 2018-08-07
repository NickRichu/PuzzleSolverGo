// SiftTest project main.go
package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	//	"sort"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

func main() {

	imageNames := readImagesIn("/home/nick/Documents/ComputerVision/hwFour/hawaii/pieces_aligned")
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

	return imageList[:5]
}

func showImage(windowName string, image gocv.Mat) {
	window := gocv.NewWindow(windowName)
	window.IMShow(image)
	gocv.WaitKey(0)
}

func puzzleSolverWithoutRotation(guideImage gocv.Mat, imageList []gocv.Mat) {
	var badImages []gocv.Mat
	var goodMatches []gocv.DMatch

	mask := gocv.NewMat()
	si := contrib.NewSIFT()
	bf := gocv.NewBFMatcher()
	kpGuide, desc1 := si.DetectAndCompute(guideImage, mask)
	for _, singleImage := range imageList {
		_, desc2 := si.DetectAndCompute(singleImage, mask)
		if !desc2.Empty() {
			match := bf.KnnMatch(desc1, desc2, 10)
			goodMatches = goodMatchFilter(match)
			for _, match := range goodMatches[:1] {
				currentKP1 := kpGuide[match.QueryIdx]
				locationx := int(math.RoundToEven(currentKP1.X)/50) * 50
				locationy := int(math.RoundToEven(currentKP1.Y)/50) * 50
				print("----------------------------------\n")
				print(locationx, " coordinate x \n")
				print(locationy, " coordinate y\n\n")

			}

		} else {
			badImages = append(badImages, singleImage)
			print("it hit this")
		}
	}

}

func goodMatchFilter(match [][]gocv.DMatch) []gocv.DMatch {
	var goodMatches []gocv.DMatch
	for i := 0; i < len(match); i++ {
		if match[i][0].Distance < 0.5*match[i][1].Distance {
			goodMatches = append(goodMatches, match[i][0])
		}
	}
	return goodMatches
}
