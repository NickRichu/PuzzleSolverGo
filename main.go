// SiftTest project main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
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
	var matches [][][]gocv.DMatch

	mask := gocv.NewMat()
	si := contrib.NewSIFT()
	bf := gocv.NewBFMatcher()

	_, desc1 := si.DetectAndCompute(guideImage, mask)

	for _, singleImage := range imageList {
		_, desc2 := si.DetectAndCompute(singleImage, mask)
		if desc2.Empty() {
			badImages = append(badImages, singleImage)
		}
		match := bf.KnnMatch(desc1, desc2, 2)
		//		for i := 0; i < len(match); i++ {
		//			if match[i][0].Distance < 0.8*match[i][1].Distance {
		//				goodMatches = append(goodMatches, match[i][0])

		//			}
		//		}
		matches = append(matches, match)

	}
	print(matches[0])

}
