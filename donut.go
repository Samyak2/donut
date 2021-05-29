package main

import (
	"fmt"
	"math"
	"time"
)

const height = 30
const width = 30
const depth = 10

const framerateX = 30
const framerateZ = 10
const framedelay = 100

const resolution = 100

// R1
const radius = 5

// R2
const offset = 5

// K_1
const cameraDist = 10

// K_2
const donutDist = 15

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

func drawScreen(zBuffer [][][2]int) {
	fmt.Print("\033[2J\033[H")
	for i := 0; i < len(zBuffer); i++ {
		for j := 0; j < len(zBuffer[i]); j++ {
			if zBuffer[i][j][0] == math.MaxInt64 {
				fmt.Print(" ")
			} else {
				fmt.Printf("%c", charMap[zBuffer[i][j][1]])
			}
		}
		fmt.Println()
	}
}

func main() {
	theta := 0.0
	phi := 0.0

	maxX := 0.0
	maxY := 0.0
	maxZ := 0.0

	minX := 0.0
	minY := 0.0
	minZ := 0.0

	zBuffer := make([][][2]int, height)
	for i := 0; i < height; i++ {
		zBuffer[i] = make([][2]int, width)
		// for j := 0; j < width; j++ {
		// 	zBuffer[i][j] = make([]float64, depth)
		// }
	}
	resetZBuffer(zBuffer)

	for a := 0; a < framerateX; a++ {
		A := (float64(a) / framerateX) * (2.0 * math.Pi)

		for b := 0; b < framerateZ; b++ {
			B := (float64(b) / framerateZ) * (2.0 * math.Pi)

			resetZBuffer(zBuffer)

			// outer loop for phi
			for i := 0; i < resolution; i++ {
				phi = (float64(i) / resolution) * (2.0 * math.Pi)
				// inner loop for theta
				for j := 0; j < resolution; j++ {
					theta = (float64(j) / resolution) * (2.0 * math.Pi)

					// pre-projection
					oldX := (math.Cos(phi)*math.Cos(B)+
						math.Sin(A)*math.Sin(B)*math.Sin(phi))*
						(offset+radius*math.Cos(theta)) -
						radius*math.Cos(A)*math.Sin(B)*math.Sin(theta)
					oldY := (offset+radius*math.Cos(theta))*
						(math.Cos(phi)*math.Sin(B)-
							math.Cos(B)*math.Sin(A)*math.Sin(phi)) +
						radius*math.Cos(A)*math.Cos(B)*math.Sin(theta)
					oldZ := (offset+radius*math.Cos(theta))*math.Cos(A)*math.Sin(phi) +
						radius*math.Sin(A)*math.Sin(theta)

					x := ((cameraDist * oldX) / (donutDist + oldZ)) + height/2
					y := ((cameraDist * oldY) / (donutDist + oldZ)) + width/2
					z := oldZ

					// fmt.Println(x, y, z)

					rX := int(math.Round(x))
					rY := int(math.Round(y))
					rZ := int(math.Round(z))

					if rZ < zBuffer[rX][rY][0] {
						// luminance
						oldL := math.Cos(phi)*math.Cos(theta)*math.Sin(B) -
							math.Cos(A)*math.Cos(theta)*math.Sin(phi) -
							math.Sin(A)*math.Sin(theta) +
							math.Cos(B)*(math.Cos(A)*math.Sin(theta)-
								math.Cos(theta)*math.Sin(A)*math.Sin(phi))

						oldL += 1.5
						oldL *= 3

						// fmt.Println(oldL)
						l := int(math.Round(oldL))
						// fmt.Println(l)

						zBuffer[rX][rY] = [2]int{rZ, l}
					}

					maxX = math.Max(maxX, x)
					maxY = math.Max(maxY, y)
					maxZ = math.Max(maxZ, z)

					minX = math.Min(minX, x)
					minY = math.Min(minY, y)
					minZ = math.Min(minZ, z)
				}
			}
			drawScreen(zBuffer)
			time.Sleep(time.Millisecond * framedelay)
		}
	}

	// fmt.Println(maxX, maxY, maxZ)
	// fmt.Println(1/maxX, 1/maxY, 1/maxZ)
	// fmt.Println(minX, minY, minZ)
	// fmt.Println(1/minX, 1/minY, 1/minZ)

}
