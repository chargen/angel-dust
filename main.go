package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"image"
	"image/color"
	"image/draw"
	"image/png"
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

func convertByteMap(byteMap [][]int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(img, img.Bounds(),
		&image.Uniform{color.Black}, image.ZP, draw.Src)

	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			if byteMap[x][y] > 0 {
				img.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
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
	img := convertByteMap(byteMap)
	saveImage(img, os.Args[2])
}
