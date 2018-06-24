package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fogleman/gg"
)

func writeToSpriteSheet(pngs []string) {
	/*
		// TODO
		for _, png := range pngs {
			fileExt := filepath.Ext(png)
			nameNoExt := strings.TrimSuffix(filepath.Base(png), fileExt)
		}
	*/

	images := map[string]image.Image{}

	var wg sync.WaitGroup
	wg.Add(len(pngs))

	for _, path := range pngs {
		go func(path string) {
			defer wg.Done()
			img, err := gg.LoadPNG(path)
			if err != nil {
				fmt.Println(err)
				return
			}
			images[path] = img
		}(path)
	}

	wg.Wait()

	imageCount := len(pngs)

	fmt.Println("Packing", imageCount, "tiles")

	imagesToConv, idx := make([]image.Image, imageCount), 0
	for _, image := range images {
		imagesToConv[idx] = image
		idx++
	}

	tileSize := imagesToConv[0].Bounds().Dx()

	imageCountHalf := imageCount / 2
	spriteSheetSize := (tileSize * imageCountHalf)
	newImage := gg.NewContext(spriteSheetSize, spriteSheetSize)

	println("packing", imageCount, "sprites into", spriteSheetSize, "x", spriteSheetSize, "sheet")

	for i, img := range imagesToConv {
		x := (i % imageCountHalf) * tileSize
		y := (i / imageCountHalf) * tileSize
		newImage.DrawImage(img, x, y)
	}

	err := newImage.SavePNG(fmt.Sprintf("spritesheet_gen_%d.png", time.Now().Unix()))
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println(`SSM - SpriteSheetMaker
Copyright (c) 2018, Felix Angell.

Note: SSM works in the preconditions that
your tiles are rectangular, i.e. S x S in size.
`)

	folder := os.Args[1]

	fmt.Println("# Searching", folder)

	absPath, err := filepath.Abs(folder)
	if err != nil {
		panic(err)
	}

	pngs := []string{}

	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".png" {
			fmt.Println("- ", path)
			pngs = append(pngs, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	start := time.Now()

	if len(pngs) > 0 {
		writeToSpriteSheet(pngs)
	}

	elapsed := time.Now().Sub(start)
	elapsedMS := elapsed.Nanoseconds() / int64(time.Millisecond)

	fmt.Println("Done!")
	fmt.Println("- Took", fmt.Sprintf("%d", elapsedMS)+"/ms")
}
