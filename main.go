package main

import (
	"fmt"         // Format
	"image"       // Image library
	"image/color" // Image sub-library for color
	"os"          // System Library

	"gocv.io/x/gocv" // GoCV Library (Wrapper for OpenCV)
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Create a OpenCV window! [filename]") // Make sure the user inputs an image file
		return
	}

	filepath := os.Args[1] // Assign image path

	window := gocv.NewWindow("Hough Transformations! | Nabeel Ahmed") // Generate window object
	defer window.Close()                                              // Window error handling

	Threshold := window.CreateTrackbar("Threshold", 25) // Window Track object
	Threshold.SetPos(100)                               // Threshold set default position

	inputImg := gocv.IMRead(filepath, gocv.IMReadGrayScale) // Input image matrix and Grayscale conversion
	defer inputImg.Close()                                  // Input image error handing

	gocv.MedianBlur(inputImg, &inputImg, 5) // Apply blur to reduce noise
	for {
		finalCompliedImage := gocv.NewMat() // Final rendered image
		defer finalCompliedImage.Close()    // Final rendered image error handling

		gocv.CvtColor(inputImg, &finalCompliedImage, gocv.ColorGrayToBGR) // Convert final image to 24 bit color map

		houghTransformation := gocv.NewMat() // Hough Transformation matrix
		defer houghTransformation.Close()    // Hough Transformation error handling

		blue := color.RGBA{0, 0, 255, 0} // blue color rgba
		red := color.RGBA{255, 0, 0, 0}  // red color rgba

		houghThreshold := Threshold.GetPos() * 100 // set hough Threshold
		gocv.HoughCirclesWithParams(               // Apply hough algorithm
			inputImg,                   // Input: input image
			&houghTransformation,       // Output: Hough Transfomation Matrix pointer
			gocv.HoughGradient,         // Hough Gradient Calculation method
			1.5,                        // dp (inverse ratio for circle accuracy)
			float64(inputImg.Rows()/8), // circle origin minium distance
			float64(houghThreshold),    // High threshold (provided by user input)
			20,                         // Low threshold
			1,                          // min radius
			0)                          // max radius (0 = unlimited radius)

		for i := 0; i < houghTransformation.Cols(); i++ { // Scan image colums for hough transformation vectors
			v := houghTransformation.GetVecfAt(0, i) // vector of transformation
			// if circles are found
			if len(v) > 2 {
				x := int(v[0]) // Image x
				y := int(v[1]) // Image y
				r := int(v[2]) // Circle radius

				gocv.Circle(&finalCompliedImage, image.Pt(x, y), r, blue, 2) //render circle
				gocv.Circle(&finalCompliedImage, image.Pt(x, y), 2, red, 3)  // render origin
			}
		}

		window.IMShow(finalCompliedImage) // Render compied image matrix to window object

		if window.WaitKey(10) >= 0 { // close window after 10 key presses (just do ctrl + c)
			break
		}
	}
}
