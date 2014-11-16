package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"sort"
)

func printUsage() {
	fmt.Println("Usage: go run main.go <input_file> <output_image>")
}

func readBytes(path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}

func mapBytes(bytes []byte) [][]int {
	byteMap := make([][]int, 256)
	for i := range byteMap {
		byteMap[i] = make([]int, 256)
	}

	for i := 0; i < len(bytes) - 1; i++ {
		byteMap[bytes[i]][bytes[i + 1]] += 1
	}

	return byteMap
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func transformByteMap(byteMap [][]int) [][]int {
	max := byteMap[0][0]
	min := byteMap[0][0]
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			max = Max(max, byteMap[x][y])
			min = Min(min, byteMap[x][y])
		}
	}

	if min == max {
		panic("min == max")
	}

	heightMap := make(map[int]int)
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			heightMap[byteMap[x][y]] = 0
		}
	}

	var heights []int
	for key := range heightMap {
		heights = append(heights, key)
	}
	sort.Ints(heights)

	step := 256.0 / float32(len(heights))
	for i, height := range heights {
		heightMap[height] = int(step * float32(i))
	}

	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			byteMap[x][y] = heightMap[byteMap[x][y]]
		}
	}

	return byteMap
}

func convertByteMap(byteMap [][]int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(img, img.Bounds(),
		&image.Uniform{color.Black}, image.ZP, draw.Src)

	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			colorValue := uint8(byteMap[x][y])
			img.Set(x, y, color.RGBA{colorValue,
				colorValue, colorValue, 255})
		}
	}

	return img
}

func saveImage(img image.Image, path string) {
	output, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	png.Encode(output, img)
}

func main() {
	if len(os.Args) != 3 {
		printUsage()
		return
	}

	bytes := readBytes(os.Args[1])
	byteMap := mapBytes(bytes)
	byteMap = transformByteMap(byteMap)
	img := convertByteMap(byteMap)
	saveImage(img, os.Args[2])
}
