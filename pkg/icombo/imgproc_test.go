package icombo

import (
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
)

var testConfig = ProcessImagesInput{
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

	testOutputDirectory := "./test_image_output"

	viper.Set("image_output_directory", testOutputDirectory)
	viper.Set("image_part_size_pixels", 32)
	viper.Set("image_input_directory", "../../example/image_parts")

	// Return a function to teardown the test
	return func(tb testing.TB) {
		log.Println("teardown suite")
		os.RemoveAll(testOutputDirectory)
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
