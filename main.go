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

func print_usage() {
	fmt.Printf("Usage: go run main.go <input_file> <output_image>\n",
		os.Args[0])
}

func main() {
	if len(os.Args) != 3 {
		print_usage()
		return
	}

	input := os.Args[1]
	output := os.Args[2]

	bytes, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}

	byte_map := make([][]int, 256)
	for i := range byte_map {
		byte_map[i] = make([]int, 256)
	}

	for i := 0; i < len(bytes) - 1; i++ {
		byte_map[bytes[i]][bytes[i + 1]] += 1
	}

	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(img, img.Bounds(),
		&image.Uniform{color.Black}, image.ZP, draw.Src)

	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			if byte_map[x][y] > 0 {
				img.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	output_file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer output_file.Close()

	png.Encode(output_file, img)
}
