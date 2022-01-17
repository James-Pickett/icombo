package icombo

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/disintegration/imaging"
)

type ProcessImagesInput struct {
	ImageDefs []ImageDef           `mapstructure:"images"`
	Options   ProcessImagesOptions `mapstructure:"options"`
}

type ProcessImagesOptions struct {
	ImageOutputDirectory string `mapstructure:"image_output_directory"`
	ImageInputDirectory  string `mapstructure:"image_input_directory"`
	Concurrency          int    `mapstructure:"concurrency"`
	ImagePartSizePixels  int    `mapstructure:"image_part_size_pixels"`
}

type ImageDef struct {
	Name          string         `mapstructure:"name"`
	ImagePartDefs []ImagePartDef `mapstructure:"image_parts"`
}

func (img ImageDef) imagePartCount() int {
	count := 0
	for i := 0; i < len(img.ImagePartDefs); i++ {
		count += img.ImagePartDefs[i].countAtLeast1()
	}
	return count
}

type ImagePartDef struct {
	FileName string  `mapstructure:"file_name"`
	Rotation float64 `mapstructure:"rotation_degrees"`
	Count    int     `mapstructure:"count"`
}

func (i ImagePartDef) countAtLeast1() int {
	count := i.Count
	if count > 0 {
		return count
	}
	return 1
}

func ProcessImages(input ProcessImagesInput) error {
	imageCount := len(input.ImageDefs)

	imageDefsChan := make(chan ImageDef, imageCount)
	errorChan := make(chan error, imageCount)

	concurrency := input.Options.Concurrency
	if input.Options.Concurrency <= 0 || input.Options.Concurrency > imageCount {
		concurrency = imageCount
	}

	// create output dir if not exists
	if _, err := os.Stat(input.Options.ImageOutputDirectory); os.IsNotExist(err) {
		os.Mkdir(input.Options.ImageOutputDirectory, os.ModeAppend)
	}

	for i := 0; i < concurrency; i++ {
		go processImageWorker(imageDefsChan, errorChan, input.Options)
	}

	for i := 0; i < imageCount; i++ {
		imageDefsChan <- input.ImageDefs[i]
	}

	var lastErr error = nil
	for i := 0; i < imageCount; i++ {
		err := <-errorChan
		if err != nil {
			lastErr = err
			log.Println(err)
		}
	}

	return lastErr
}

func processImageWorker(imagedDefs <-chan ImageDef, errors chan<- error, opts ProcessImagesOptions) {
	for imageDef := range imagedDefs {
		errors <- createImage(imageDef, opts)
	}
}

func createImage(imageDef ImageDef, opts ProcessImagesOptions) error {

	finalImage := imaging.New(opts.ImagePartSizePixels*imageDef.imagePartCount(), opts.ImagePartSizePixels, color.Black)
	totalImagePartCount := 0

	for i := 0; i < len(imageDef.ImagePartDefs); i++ {

		imagePartDef := imageDef.ImagePartDefs[i]

		for j := 0; j < imagePartDef.countAtLeast1(); j++ {

			imagePartPath := fmt.Sprint(opts.ImageInputDirectory, "/", resolvePngExtention(imagePartDef.FileName))

			rawImagePartSrc, err := imaging.Open(imagePartPath)
			if err != nil {
				return err
			}

			imagePartSrc := imaging.Resize(rawImagePartSrc, opts.ImagePartSizePixels, opts.ImagePartSizePixels, imaging.Lanczos)
			imagePartSrc = imaging.Rotate(imagePartSrc, imagePartDef.Rotation, color.Black)
			finalImage = imaging.Paste(finalImage, imagePartSrc, image.Point{opts.ImagePartSizePixels * totalImagePartCount, 0})

			totalImagePartCount++
		}
	}

	fileOutputPath := fmt.Sprint(opts.ImageOutputDirectory, "/", resolvePngExtention(imageDef.Name))

	if err := saveImage(finalImage, fileOutputPath); err != nil {
		return err
	}

	return nil
}

func saveImage(img *image.NRGBA, path string) error {
	if err := imaging.Save(img, path); err != nil {
		return err
	}
	return nil
}

func resolvePngExtention(path string) string {
	if strings.HasSuffix(path, ".png") {
		return path
	}
	return fmt.Sprint(path, ".png")
}
