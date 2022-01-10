package icombo

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/spf13/viper"
)

type ProcessImagesInput struct {
	ImageDefs []ImageDef `mapstructure:"images"`
}

type ImageDef struct {
	Name          string         `mapstructure:"name"`
	ImagePartDefs []ImagePartDef `mapstructure:"image_parts"`
}

func (i ImageDef) fileOutputPath() string {
	outputDir := viper.Get("image_output_directory")
	path := fmt.Sprint(outputDir, "/", i.Name)
	return resolvePngExtention(path)
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

func (i ImagePartDef) filePath() string {
	inputDir := viper.Get("image_input_directory")
	path := fmt.Sprint(inputDir, "/", i.FileName)
	return resolvePngExtention(path)
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

	concurrency := viper.GetInt("concurrency")
	if concurrency <= 0 || concurrency > imageCount {
		concurrency = imageCount
	}

	for i := 0; i < concurrency; i++ {
		go processImageWorker(imageDefsChan, errorChan)
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

func processImageWorker(imagedDefs <-chan ImageDef, errors chan<- error) {
	for imageDef := range imagedDefs {
		errors <- createImage(imageDef)
	}
}

func createImage(imageDef ImageDef) error {

	finalImage := imaging.New(imagePartSize()*imageDef.imagePartCount(), imagePartSize(), color.Black)
	totalImagePartCount := 0

	for i := 0; i < len(imageDef.ImagePartDefs); i++ {

		imagePartDef := imageDef.ImagePartDefs[i]

		for j := 0; j < imagePartDef.countAtLeast1(); j++ {
			rawImagePartSrc, err := imaging.Open(imagePartDef.filePath())
			if err != nil {
				return err
			}

			imagePartSrc := imaging.Resize(rawImagePartSrc, imagePartSize(), imagePartSize(), imaging.Lanczos)
			imagePartSrc = imaging.Rotate(imagePartSrc, imagePartDef.Rotation, color.Black)
			finalImage = imaging.Paste(finalImage, imagePartSrc, image.Point{imagePartSize() * totalImagePartCount, 0})

			totalImagePartCount++
		}
	}

	if err := saveImage(finalImage, imageDef.fileOutputPath()); err != nil {
		return err
	}

	return nil
}

func saveImage(img *image.NRGBA, path string) error {
	outputDir := viper.GetString("image_output_directory")
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModeAppend)
	}
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

func imagePartSize() int {
	return viper.GetInt("image_part_size_pixels")
}
