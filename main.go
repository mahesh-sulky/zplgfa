package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"simonwaldherr.de/go/zplgfa/zplgfa"
)

func main() {
	startTime := time.Now()

	// open file
	file, err := os.Open("input.jpg")
	if err != nil {
		log.Printf("Warning: could not open the file: %s\n", err)
		return
	}

	defer file.Close()

	// load image head information
	config, format, err := image.DecodeConfig(file)
	if err != nil {
		log.Printf("Warning: image not compatible, format: %s, config: %v, error: %s\n", format, config, err)
	}

	// reset file pointer to the beginning of the file
	file.Seek(0, 0)

	// load and decode image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Printf("Warning: could not decode the file, %s\n", err)
		return
	}

	// Rotate the image 90 degrees to the right
	rotatedImg := imaging.Rotate270(img)

	// flatten image
	flat := zplgfa.FlattenImage(rotatedImg)

	// convert image to zpl compatible type
	gfimg := zplgfa.ConvertToZPL(flat, zplgfa.CompressedASCII)

	// create and open a .zpl file for writing
	outputFile, err := os.Create("output.zpl")
	if err != nil {
		log.Printf("Warning: could not create the output file: %s\n", err)
		return
	}
	defer outputFile.Close()

	// write the ZPL data to the file
	_, err = outputFile.WriteString(gfimg)
	if err != nil {
		log.Printf("Warning: could not write to the output file: %s\n", err)
		return
	}

	// Calculate and display elapsed time
	elapsedTime := time.Since(startTime)
	fmt.Printf("Processing time: %s\n", elapsedTime)
}
