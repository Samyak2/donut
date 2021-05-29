package main

import (
	"fmt"
	"math"
)

const height = 15
const width = 15
const depth = 10

// R1
const radius = 5

// R2
const offset = 5

// K_1
const cameraDist = 10

// K_2
const donutDist = 20

func drawScreen(zBuffer [][]int) {
	for i := 0; i < len(zBuffer); i++ {
		for j := 0; j < len(zBuffer[i]); j++ {
			if zBuffer[i][j] == math.MaxInt64 {
				fmt.Print(" ")
			} else {
				fmt.Print("•")
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

	zBuffer := make([][]int, height)
	for i := 0; i < height; i++ {
		zBuffer[i] = make([]int, width)

		for j := 0; j < width; j++ {
			zBuffer[i][j] = math.MaxInt64
		}

		// for j := 0; j < width; j++ {
		// 	zBuffer[i][j] = make([]float64, depth)
		// }
	}

	// outer loop for phi
	// inner loop for theta
	for i := 0; i < height; i++ {
		phi = (float64(i) / height) * (2.0 * math.Pi)
		for j := 0; j < width; j++ {
			theta = (float64(j) / width) * (2.0 * math.Pi)

			// pre-projection
			oldX := math.Cos(phi) * (offset + radius * math.Cos(theta))
			oldY := radius * math.Sin(theta)
			oldZ := -math.Sin(phi) * (offset + radius * math.Cos(theta))

			x := ((cameraDist * oldX) / (donutDist + oldZ)) + height / 2
			y := ((cameraDist * oldY) / (donutDist + oldZ)) + width / 2
			z := oldZ

			// fmt.Println(x, y, z)

			rX := int(math.Round(x))
			rY := int(math.Round(y))
			rZ := int(math.Round(z))

			zBuffer[rX][rY] = rZ

			maxX = math.Max(maxX, x)
			maxY = math.Max(maxY, y)
			maxZ = math.Max(maxZ, z)

			minX = math.Min(minX, x)
			minY = math.Min(minY, y)
			minZ = math.Min(minZ, z)
		}
	}

	// fmt.Println(maxX, maxY, maxZ)
	// fmt.Println(1/maxX, 1/maxY, 1/maxZ)
	// fmt.Println(minX, minY, minZ)
	// fmt.Println(1/minX, 1/minY, 1/minZ)

	drawScreen(zBuffer)
}
