package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/nfnt/resize"
)

// https://github.com/nfnt/resize
func CompressImage(inputImage *bytes.Buffer, maxDimension uint) (*bytes.Buffer, string, error) {
	// Decode the image
	img, format, err := image.Decode(inputImage)
	if err != nil {
		return nil, "", err
	}

	// Calculate the aspect ratio
	aspectRatio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())

	// Calculate the new dimensions based on the maximum dimension
	var newWidth, newHeight uint
	if img.Bounds().Dx() > img.Bounds().Dy() {
		newWidth = maxDimension
		newHeight = uint(float64(maxDimension) / aspectRatio)
	} else {
		newWidth = uint(float64(maxDimension) * aspectRatio)
		newHeight = maxDimension
	}

	// Resize the image while maintaining the aspect ratio
	resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	// Create a buffer to store the compressed image
	compressedBuffer := new(bytes.Buffer)

	// Encode the resized image to the buffer based on the image format
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(compressedBuffer, resizedImg, nil)
	case "png":
		err = png.Encode(compressedBuffer, resizedImg)
	case "gif":
		err = gif.Encode(compressedBuffer, resizedImg, nil)
	default:
		return nil, "", fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return nil, "", err
	}

	return compressedBuffer, fmt.Sprintf(".%s", format), nil
}

// import (
// 	"bytes"
// 	"fmt"
// 	"image"
// 	_ "image/jpeg" // Register JPEG format
// 	_ "image/png"  // Register PNG format
// 	"io/ioutil"

// 	"github.com/rwcarlsen/goexif/exif"
// )

// func extractImageMetadata(imageBytes []byte, filename string) (map[string]interface{}, error) {
// 	// Size
// 	size := len(imageBytes)

// 	// Dimensions
// 	img, _, err := image.Decode(bytes.NewReader(imageBytes))
// 	if err != nil {
// 		return nil, fmt.Errorf("error decoding image: %w", err)
// 	}
// 	bounds := img.Bounds()
// 	width, height := bounds.Max.X, bounds.Max.Y

// 	// EXIF Data
// 	exifData, err := exif.Decode(bytes.NewReader(imageBytes))
// 	if err != nil {
// 		return nil, fmt.Errorf("error decoding EXIF data: %w", err)
// 	}

// 	// Extracting EXIF fields (example: DateTime and CameraModel)
// 	dateTime, _ := exifData.DateTime()
// 	cameraModel, _ := exifData.Get(exif.Model)

// 	// Compile the metadata
// 	metadata := map[string]interface{}{
// 		"Filename":     filename,
// 		"Size":         size,
// 		"Width":        width,
// 		"Height":       height,
// 		"DateTime":     dateTime,
// 		"CameraModel":  cameraModel.StringVal(),
// 	}

// 	return metadata, nil
// }

// func main() {
// 	// Example usage
// 	imageBytes, err := ioutil.ReadFile("example.jpg") // Replace with your image source
// 	if err != nil {
// 		panic(err)
// 	}
// 	filename := "example.jpg" // Replace with the actual filename if available

// 	metadata, err := extractImageMetadata(imageBytes, filename)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("Metadata: %+v\n", metadata)
// }
