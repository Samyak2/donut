package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// some basic configuration
// TODO: expose these as flags
const size = 40
const aspectRatio = 1.75

// select colors here
var charMap = charMapBW
// set only when charMap is charMapBW (otherwise 0)
var brightness = 2

// these are actually interchanged lol
const height = size * aspectRatio
const width = size

// start angle in X
const startX = 20.0 * math.Pi / 180.0
// angle to move in every step in X direction
// determine smoothness and speed of rotation
const stepX = 1.0 * math.Pi / 180.0

// start angle in Z
const startZ = 30.0 * math.Pi / 180.0
// similarly for Z
// 360 (or more) means it never rotates this way
const stepZ = 360.0 * math.Pi / 180.0

// frametime in milliseconds
// 16 ~> 60fps
// 32 ~> 30fps
const framedelay = 32

// lower res can cause black spots to appear
// 6 was decided arbitrarily
const resolutionPhi = size * 6
const resolutionTheta = size * 6

// R1
const radius = 1.0

// R2
const offset = 2.0

// K_2
const donutDist = 5.0

// K_1
const cameraDist = width * donutDist * 3 / (8 * (radius + offset))
// K_1 in y set differently due to terminal text caret dimensions
//   being rectangular
const cameraDistY = height * donutDist * 3 / (8 * (radius + offset))

// colors!
// refer to https://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html#256-colors
//   for the numbers
// goes from dark to light
var charMapPink = []int{196, 197, 198, 199, 200, 201, 205, 206, 207, 218, 219, 224, 225}
var charMapBW = []int{233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244}
var charMapBlue = []int{16, 17, 18, 19, 20, 21, 24, 25, 26, 27, 32, 33, 69}


// sets z-buffer to Inf (z-value) and 0 (luminance)
func resetZBuffer(zBuffer [][][2]int) {
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			zBuffer[i][j] = [2]int{math.MaxInt64, 0}
		}
	}
}

// draw z-buffer on stdout
// a bufio.Writer is used to prevent buffering on stdout
//   this removes flickering
func drawScreen(f *bufio.Writer, zBuffer [][][2]int) {
	// flush buffer AFTER writing
	//  prevents tearing/flickering
	defer f.Flush()

	// this moves cursor to start of the screen
	f.WriteString("\033[H")

	for i := 0; i < len(zBuffer); i++ {
		for j := 0; j < len(zBuffer[i]); j++ {
			if zBuffer[i][j][0] == math.MaxInt64 {
				// nothing to render at this point
				f.WriteString(" ")
			} else {
				// draw the full block character █
				//  along with the color
				f.WriteString(fmt.Sprintf("\u001b[38;5;%dm█", charMap[zBuffer[i][j][1]] + brightness))
			}
		}
		f.WriteString("\n")
	}
}

func main() {
	// handle ctrl-c gracefully
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// reset cursor back to normal
		fmt.Print("\033[?25h")
		os.Exit(1)
	}()
	//

	theta := 0.0
	phi := 0.0

	// the z-buffer consists of 2 values for every point
	//   - z-value: value in z dimension
	//   - luminance: for shading (refer article for details)
	zBuffer := make([][][2]int, width)
	for i := 0; i < width; i++ {
		zBuffer[i] = make([][2]int, height)
	}
	resetZBuffer(zBuffer)

	f := bufio.NewWriter(os.Stdout)

	// clear screen and hide cursor
	fmt.Print("\033[2J\033[?25l")

	// A is angle in X dimension
	A := startX
	// B is angle in Z dimension
	B := startZ

	for ; A <= 2.0*math.Pi+stepX; A += stepX {
		// pre-compute these
		cosA := math.Cos(A)
		sinA := math.Sin(A)

		for ; B <= 2.0*math.Pi+stepZ; B += stepZ {
			cosB := math.Cos(B)
			sinB := math.Sin(B)

			// ready for rendering
			resetZBuffer(zBuffer)

			// outer loop for phi
			for i := 0; i < resolutionPhi; i++ {
				phi = (float64(i) / resolutionPhi) * (2.0 * math.Pi)
				cosΦ := math.Cos(phi)
				sinΦ := math.Sin(phi)

				// inner loop for theta
				for j := 0; j < resolutionTheta; j++ {
					theta = (float64(j) / resolutionTheta) * (2.0 * math.Pi)
					cosθ := math.Cos(theta)
					sinθ := math.Sin(theta)

					circleX := offset + radius*cosθ
					circleY := radius * sinθ

					// pre-projection
					// refer article for the math
					oldX := (cosΦ*cosB+sinA*sinB*sinΦ)*circleX -
						circleY*cosA*sinB
					oldY := circleX*(cosΦ*sinB-cosB*sinA*sinΦ) +
						circleY*cosA*cosB
					oldZ := donutDist + (circleX)*cosA*sinΦ + circleY*sinA

					// project onto the screen
					// add width/2 since our screen starts at top left
					//   (as opposed to the middle in cartesian coords)
					x := width/2 + ((cameraDist * oldX) / (oldZ))
					// similarly for height, but y-axis goes downwards
					//   (as opposed to upwards for +ve in cartesian)
					y := height/2 - ((cameraDistY * oldY) / (oldZ))
					// note: unlike the article, I use z directly instead of
					//    1/z
					z := oldZ

					// fmt.Println(x, y, z)

					// discretize
					// math.Round here does not bring much improvement
					rX := int(x)
					rY := int(y)
					rZ := int(z)

					// if current point is closer than the one already
					// in the z-buffer, we overwrite it
					if rZ < zBuffer[rX][rY][0] {
						// luminance
						oldL := cosΦ*cosθ*sinB -
							cosA*cosθ*sinΦ -
							sinA*sinθ +
							cosB*(cosA*sinθ-
								cosθ*sinA*sinΦ)

						// we only care when the luminance is positive
						if oldL > 0 {
							// luminance ranges from -sqrt(2) to +sqrt(2)
							// scale it to 0 to 11.3 and discretize
							l := int(oldL * 8)

							zBuffer[rX][rY] = [2]int{rZ, l}
						} else {
							zBuffer[rX][rY] = [2]int{rZ, 0}
						}
					}
				}
			}
			drawScreen(f, zBuffer)
			// maintain framerate
			time.Sleep(time.Millisecond * framedelay)
		}
		// prevent this going to infinity, just in case
		B -= 2.0 * math.Pi
	}
	A -= 2.0 * math.Pi

	// reset cursor back to normal
	fmt.Print("\033[?25h")
}
