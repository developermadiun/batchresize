package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
)

func main() {
	// Open the list.txt file
	fileList, err := os.Open("list.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fileList.Close()

	// Create a scanner to read the file names line by line
	scanner := bufio.NewScanner(fileList)

	// Create the result folder if it doesn't exist
	resultFolder := "result"
	if _, err := os.Stat(resultFolder); os.IsNotExist(err) {
		err := os.Mkdir(resultFolder, os.ModeDir|os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Process each line in the file
	for scanner.Scan() {
		fileName := scanner.Text()

		// Combine the file name with the uploads and result folder path
		uploadFilePath := filepath.Join("uploads", fileName)
		resultFilePath := filepath.Join(resultFolder, fileName)

		// Check if the file exists
		if _, err := os.Stat(uploadFilePath); err == nil {
			fmt.Printf("Processing file: %s\n", fileName)

			// Open the image file
			imgFile, err := os.Open(uploadFilePath)
			if err != nil {
				log.Println("Error opening image file:", err)
				continue
			}
			defer imgFile.Close()

			// Decode the WebP image
			img, err := webp.Decode(imgFile)
			if err != nil {
				log.Println("Error decoding WebP image:", err)
				continue
			}

			// Crop the image to a square in the center
			croppedImg := cropToSquare(img)

			// Save the cropped image to the result folder
			resultFile, err := os.Create(resultFilePath)
			if err != nil {
				log.Println("Error creating result file:", err)
				continue
			}
			defer resultFile.Close()

			// Encode the cropped image to WebP format
			err = webp.Encode(resultFile, croppedImg, &webp.Options{Quality: 100})
			if err != nil {
				log.Println("Error encoding WebP image:", err)
				continue
			}

			fmt.Printf("Successfully processed: %s\n", fileName)
		} else {
			fmt.Printf("File not found: %s\n", fileName)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// cropToSquare crops the image to a square in the center
func cropToSquare(img image.Image) image.Image {
	bounds := img.Bounds()
	minDim := bounds.Dx()
	if bounds.Dy() < minDim {
		minDim = bounds.Dy()
	}

	// Calculate the center coordinates
	centerX := bounds.Min.X + bounds.Dx()/2
	centerY := bounds.Min.Y + bounds.Dy()/2

	// Calculate the cropping boundaries
	cropMinX := centerX - minDim/2
	cropMinY := centerY - minDim/2
	cropMaxX := centerX + minDim/2
	cropMaxY := centerY + minDim/2

	croppedImg := image.NewRGBA(image.Rect(0, 0, minDim, minDim))
	for y := cropMinY; y < cropMaxY; y++ {
		for x := cropMinX; x < cropMaxX; x++ {
			croppedImg.Set(x-cropMinX, y-cropMinY, img.At(x, y))
		}
	}

	return croppedImg
}
