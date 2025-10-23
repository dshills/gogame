// Package main generates PNG assets for the Space Battle game
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	log.Println("Generating Space Battle assets...")

	// Generate player ship (blue triangle pointing up)
	if err := generatePlayerShip("assets/player.png"); err != nil {
		log.Fatal(err)
	}
	log.Println("✓ Generated player.png")

	// Generate enemy ship (red triangle pointing down)
	if err := generateEnemyShip("assets/enemy.png"); err != nil {
		log.Fatal(err)
	}
	log.Println("✓ Generated enemy.png")

	// Generate bullet (yellow rectangle)
	if err := generateBullet("assets/bullet.png"); err != nil {
		log.Fatal(err)
	}
	log.Println("✓ Generated bullet.png")

	// Generate star (white dot)
	if err := generateStar("assets/star.png"); err != nil {
		log.Fatal(err)
	}
	log.Println("✓ Generated star.png")

	log.Println("All assets generated successfully!")
}

// generatePlayerShip creates a blue triangle pointing upward
func generatePlayerShip(filename string) error {
	size := 32
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Blue color
	blue := color.RGBA{R: 100, G: 150, B: 255, A: 255}

	// Draw triangle pointing up (simple filled triangle)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Triangle math: top point at (16, 4), base from (4, 28) to (28, 28)
			cx := size / 2
			topY := 4
			baseY := 28

			// Check if point is inside triangle
			if y >= topY && y <= baseY {
				progress := float64(y-topY) / float64(baseY-topY)
				halfWidth := int(progress * float64(size/2-4))
				if x >= cx-halfWidth && x <= cx+halfWidth {
					img.Set(x, y, blue)
				}
			}
		}
	}

	return savePNG(img, filename)
}

// generateEnemyShip creates a red triangle pointing downward
func generateEnemyShip(filename string) error {
	size := 32
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Red color
	red := color.RGBA{R: 255, G: 100, B: 100, A: 255}

	// Draw triangle pointing down
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			cx := size / 2
			topY := 4
			baseY := 28

			// Inverted triangle
			if y >= topY && y <= baseY {
				progress := float64(baseY-y) / float64(baseY-topY)
				halfWidth := int(progress * float64(size/2-4))
				if x >= cx-halfWidth && x <= cx+halfWidth {
					img.Set(x, y, red)
				}
			}
		}
	}

	return savePNG(img, filename)
}

// generateBullet creates a yellow rectangle
func generateBullet(filename string) error {
	width := 8
	height := 16
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Yellow color
	yellow := color.RGBA{R: 255, G: 255, B: 100, A: 255}

	// Fill rectangle
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, yellow)
		}
	}

	return savePNG(img, filename)
}

// generateStar creates a white dot
func generateStar(filename string) error {
	size := 4
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// White color
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Fill small square
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, white)
		}
	}

	return savePNG(img, filename)
}

// savePNG saves an image to a PNG file
func savePNG(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
