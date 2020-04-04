/*
 * Insert an image to a PDF file compressed using JBIG2 encoder.
 *
 * Example go run jbig2_compressed_image_in_pdf.go /tmp/output.pdf 1 /tmp/image.jpg /tmp/image2.jpg
 * adds the image to the upper left corner of the page (0,0).  The width is 100 (typical page width 612 with defaults).
 *
 * Syntax: go run jbig2_compressed_image.go output.pdf img1.jpg, img2.jpg
 */

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run jbig2_compressed_image_in_pdf.go output.pdf img1.jpg img2.jpg ...\n")
		os.Exit(1)
	}

	outputPath := os.Args[1]
	inputPaths := os.Args[2:len(os.Args)]

	err := imagesToJBIG2ToPdf(inputPaths, outputPath)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

// Images to PDF.
func imagesToJBIG2ToPdf(inputPaths []string, outputPath string) error {
	c := creator.New()

	for _, imgPath := range inputPaths {
		common.Log.Debugf("Encoding image: %s", imgPath)
		img, err := c.NewImageFromFile(imgPath)
		if err != nil {
			common.Log.Debug("Error loading image: %v", err)
			return err
		}
		// Convert the image into binary format. The RGB and GrayScale images would be converted into bi-level image.
		// This step is required for the JBIG2 Encoder.
		if err = img.ToBinaryImage(); err != nil {
			return err
		}
		// Set the JBIG2 Encoder as the image encoder.
		e := core.NewJBIG2Encoder()
		// For images that might equal following lines it might be convenient
		// to set DuplicatedLinesRemoval option to true.
		e.DefaultPageSettings.DuplicatedLinesRemoval = true
		img.SetEncoder(e)

		img.ScaleToWidth(612.0)
		// Use page width of 612 points, and calculate the height proportionally based on the image.
		// Standard PPI is 72 points per inch, thus a width of 8.5".
		height := 612.0 * img.Height() / img.Width()
		c.SetPageSize(creator.PageSize{612, height})
		c.NewPage()
		img.SetPos(0, 0)
		if err = c.Draw(img); err != nil {
			return err
		}
	}

	err := c.WriteToFile(outputPath)
	return err
}