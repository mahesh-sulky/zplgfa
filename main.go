package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"simonwaldherr.de/go/zplgfa/zplgfa"
)

// logError logs the error with a message
func logError(err error, msg string) {
	if err != nil {
		log.Printf("Warning: %s: %s\n", msg, err)
	}
}

func main() {
	// Define command-line flags
	inputFile := flag.String("input", "", "Path to the input image file (mandatory)")
	outputFile := flag.String("output", "", "Path to the output ZPL file (mandatory)")
	inputDpi := flag.String("inputdpi", "200", "Input DPI (optional, default is 200)")
	outputDpi := flag.String("outputdpi", "300", "Output DPI (optional, default is 300)")

	// Parse command-line flags
	flag.Parse()

	// Check mandatory flags
	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Error: Input and output file paths are mandatory.")
		flag.Usage()
		return
	}

	outDPI, err := strconv.ParseFloat(*outputDpi, 64)
	if err != nil {
		log.Fatalf("Invalid DPI value: %s. Must be a number.", *outputDpi)
	}

	if outDPI <= 0 {
		log.Fatalf("Invalid DPI value: %s. Must be greater than zero.", *outputDpi)
	}

	// Parse the DPI argument
	inDPI, err := strconv.ParseFloat(*inputDpi, 64)
	if err != nil {
		log.Fatalf("Invalid DPI value: %s. Must be a number.", *inputDpi)
	}

	if inDPI <= 0 {
		log.Fatalf("Invalid DPI value: %s. Must be greater than zero.", *inputDpi)
	}

	// Open the input file
	file, err := os.Open(*inputFile)
	logError(err, "could not open the file")
	if err != nil {
		return
	}
	defer file.Close()

	// Create a buffered reader
	reader := bufio.NewReader(file)

	// Load image head information
	config, _, err := image.DecodeConfig(reader)
	logError(err, "image not compatible")
	if err != nil {
		return
	}

	// Reset file pointer to the beginning of the file
	if _, err := file.Seek(0, 0); err != nil {
		logError(err, "could not seek to the beginning of the file")
		return
	}

	// Load and decode image
	img, _, err := image.Decode(reader)
	logError(err, "could not decode the file")
	if err != nil {
		return
	}

	// Rotate the image 90 degrees to the right
	rotatedImg := imaging.Rotate270(img)

	// Flatten and resize image in one go
	flat := zplgfa.FlattenImage(rotatedImg)
	scaleFactor := outDPI / inDPI
	newHeight := int(float64(config.Width) * scaleFactor)
	newWidth := int(float64(config.Height) * scaleFactor)
	resizedImage := imaging.Resize(flat, newWidth, newHeight, imaging.Lanczos)

	// Convert image to ZPL compatible type
	gfimg := zplgfa.ConvertToZPL(resizedImage, zplgfa.CompressedASCII)

	// Create and open the output file for writing
	output, err := os.Create(*outputFile)
	logError(err, "could not create the output file")
	if err != nil {
		return
	}
	defer output.Close()

	// Create a buffered writer
	writer := bufio.NewWriter(output)

	// Write the ZPL data to the file using the buffered writer
	if _, err := writer.WriteString(gfimg); err != nil {
		logError(err, "could not write to the output file")
		return
	}

	// Flush the buffer to ensure all data is written to the file
	if err := writer.Flush(); err != nil {
		logError(err, "could not flush the buffered writer")
		return
	}
	fmt.Println("Image processing completed successfully.")

}
