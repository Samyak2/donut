package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"
)

const height = 40
const width = 40
const depth = 10

// const framerateX = 30
const startX = 0.0 * math.Pi / 180.0
const stepX = 2.0 * math.Pi / 180.0
// const framerateZ = 10
const startZ = 60.0 * math.Pi / 180.0
const stepZ = 360 * math.Pi / 180.0

const framedelay = 16

const resolutionPhi = 180
const resolutionTheta = 90

// R1
const radius = 1.0

// R2
const offset = 2.0

// K_2
const donutDist = 5.0

// K_1
const cameraDist = width * donutDist * 3 / (8 * (radius + offset))

const charMap = ".,-~:;=!*#$@"

// angle X
// const A = 0.5

// angle Z
// const B = 0.5

func resetZBuffer(zBuffer [][][2]int) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			zBuffer[i][j] = [2]int{math.MaxInt64, 0}
		}
	}
}

func drawScreen(f *bufio.Writer, zBuffer [][][2]int) {
	defer f.Flush()
	f.WriteString("\033[2J\033[H")
	for i := 0; i < len(zBuffer); i++ {
		for j := 0; j < len(zBuffer[i]); j++ {
			if zBuffer[i][j][0] == math.MaxInt64 {
				f.WriteString(" ")
			} else {
				f.WriteString(fmt.Sprintf("%c", charMap[zBuffer[i][j][1]]))
			}
		}
		f.WriteString("\n")
	}
}

func main() {
	theta := 0.0
	phi := 0.0

	// maxX := 0.0
	// maxY := 0.0
	// maxZ := 0.0

	// minX := 0.0
	// minY := 0.0
	// minZ := 0.0

	zBuffer := make([][][2]int, height)
	for i := 0; i < height; i++ {
		zBuffer[i] = make([][2]int, width)
		// for j := 0; j < width; j++ {
		// 	zBuffer[i][j] = make([]float64, depth)
		// }
	}
	resetZBuffer(zBuffer)

	f := bufio.NewWriter(os.Stdout)

	fmt.Print("\033[?25l")

	// fmt.Println(cameraDist)

	for A := startX; A < 2.0 * math.Pi; A += stepX {
		// A := (float64(a) / framerateX) * (1.0 * math.Pi)
		cosA := math.Cos(A)
		sinA := math.Sin(A)

		for B := startZ; B < 2.0 * math.Pi; B += stepZ {
			// B := (float64(b) / framerateZ) * (2.0 * math.Pi)
			cosB := math.Cos(B)
			sinB := math.Sin(B)

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
					oldX := (cosΦ*cosB+sinA*sinB*sinΦ)*circleX -
						circleY*cosA*sinB
					oldY := circleX*(cosΦ*sinB-cosB*sinA*sinΦ) +
						circleY*cosA*cosB
					oldZ := donutDist + (circleX)*cosA*sinΦ + circleY*sinA

					x := width/2 + ((cameraDist * oldX) / (oldZ))
					y := height/2 - ((cameraDist * oldY) / (oldZ))
					z := oldZ

					// fmt.Println(x, y, z)

					rX := int(math.Round(x))
					rY := int(math.Round(y))
					rZ := int(math.Round(z))

					if rZ < zBuffer[rX][rY][0] {
						// luminance
						oldL := cosΦ*cosθ*sinB -
							cosA*cosθ*sinΦ -
							sinA*sinθ +
							cosB*(cosA*sinθ-
								cosθ*sinA*sinΦ)

						if oldL > 0 {
							// fmt.Println(oldL)
							l := int(math.Round(oldL * 8))
							// fmt.Println(l)

							zBuffer[rX][rY] = [2]int{rZ, l}
						} else {
							zBuffer[rX][rY] = [2]int{rZ, 0}
						}
					}

					// maxX = math.Max(maxX, x)
					// maxY = math.Max(maxY, y)
					// maxZ = math.Max(maxZ, z)

					// minX = math.Min(minX, x)
					// minY = math.Min(minY, y)
					// minZ = math.Min(minZ, z)
				}
			}
			drawScreen(f, zBuffer)
			time.Sleep(time.Millisecond * framedelay)
		}
	}

	// fmt.Println(maxX, maxY, maxZ)
	// fmt.Println(1/maxX, 1/maxY, 1/maxZ)
	// fmt.Println(minX, minY, minZ)
	// fmt.Println(1/minX, 1/minY, 1/minZ)
	fmt.Print("\033[?25h")

}
