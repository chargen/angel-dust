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

func countFreq(bytes []byte) [][]int {
	freq := make([][]int, 256)
	for i := range freq {
		freq[i] = make([]int, 256)
	}

	for i := 0; i < len(bytes) - 1; i++ {
		freq[bytes[i]][bytes[i + 1]] += 1
	}

	return freq
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

func transformFreq(freq [][]int) [][]int {
	max := freq[0][0]
	min := freq[0][0]
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			max = Max(max, freq[x][y])
			min = Min(min, freq[x][y])
		}
	}

	if min == max {
		panic("min == max")
	}

	colorMap := make(map[int]int)
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			colorMap[freq[x][y]] = 0
		}
	}

	var heights []int
	for key := range colorMap {
		heights = append(heights, key)
	}
	sort.Ints(heights)

	step := 256.0 / float32(len(heights))
	for i, height := range heights {
		colorMap[height] = int(step * float32(i))
	}

	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			freq[x][y] = colorMap[freq[x][y]]
		}
	}

	return freq
}

func convertFreq(freq [][]int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(img, img.Bounds(),
		&image.Uniform{color.Black}, image.ZP, draw.Src)

	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			colorValue := uint8(freq[x][y])
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
	freq := countFreq(bytes)
	freq = transformFreq(freq)
	img := convertFreq(freq)
	saveImage(img, os.Args[2])
}
