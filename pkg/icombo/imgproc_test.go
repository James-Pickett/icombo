package icombo

import (
	"log"
	"os"
	"testing"
)

const TEST_IMAGE_OUTPUT_DIR = "./test_image_output"
const TEST_IMAGE_INPUT_DIR = "../../example/image_parts"
const TEST_IMAGE_PART_PIXEL_SIZE = 32

var testConfig = ProcessImagesInput{
	Options: ProcessImagesOptions{
		Concurrency:          0,
		ImageOutputDirectory: TEST_IMAGE_OUTPUT_DIR,
		ImageInputDirectory:  TEST_IMAGE_INPUT_DIR,
		ImagePartSizePixels:  TEST_IMAGE_PART_PIXEL_SIZE,
	},
	ImageDefs: []ImageDef{
		{
			Name: "boards_to_logs",
			ImagePartDefs: []ImagePartDef{
				{
					FileName: "log",
					Count:    5,
				},
				{
					FileName: "arrow.png",
					Rotation: 270,
				},
				{
					FileName: "board",
					Count:    10,
				},
			},
		}, {
			Name: "boards_to_houses",
			ImagePartDefs: []ImagePartDef{
				{
					FileName: "board",
					Count:    10,
				},
				{
					FileName: "arrow.png",
					Rotation: 270,
				},
				{
					FileName: "house",
					Count:    2,
				},
			},
		},
	},
}

func setupProcessImagesSuccess(tb testing.TB) func(tb testing.TB) {
	log.Println("setup suite")
	// Return a function to teardown the test
	return func(tb testing.TB) {
		log.Println("teardown suite")
		os.RemoveAll(TEST_IMAGE_OUTPUT_DIR)
	}
}

func TestProcessImagesSuccess(t *testing.T) {
	tearDownTest := setupProcessImagesSuccess(t)
	defer tearDownTest(t)

	t.Run("process_image_test", func(t *testing.T) {
		if err := ProcessImages(testConfig); err != nil {
			t.Error("did not expect error, got", err)
		}
	})
}
