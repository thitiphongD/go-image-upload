package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
)

func main() {
	router := fiber.New()

	router.Use(cors.New())

	router.Static("/images", "./images")

	router.Post("/", handleFileUpload)

	router.Delete("/:imageName", handleDeleteImage)

	router.Listen(":4000")
}

func handleFileUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("image")

	if err != nil {
		log.Println("upload error", err)
		return c.JSON(fiber.Map{
			"status":  500,
			"message": "Server error",
			"data":    nil,
		})
	}

	uniqueId := uuid.New()

	filename := strings.Replace(uniqueId.String(), "-", "", -1)

	fileExt := strings.Split(file.Filename, ".")[1]

	image := fmt.Sprintf("%s.%s", filename, fileExt)

	err = c.SaveFile(file, fmt.Sprintf("./images/%s", image))

	if err != nil {
		log.Println("upload error", err)
		return c.JSON(fiber.Map{
			"status":  500,
			"message": "Server error",
			"data":    nil,
		})
	}

	imageUrl := fmt.Sprintf("http://localhost:4000/images/%s", image)

	data := map[string]interface{}{
		"imageName": image,
		"imageUrl":  imageUrl,
		"header":    file.Header,
		"size":      file.Size,
	}

	return c.JSON(fiber.Map{
		"status":  201,
		"message": "Upload Successfully",
		"data":    data,
	})
}

func handleDeleteImage(c *fiber.Ctx) error {
	imageName := c.Params("imageName")

	err := os.Remove(fmt.Sprintf("./images/%s", imageName))

	if err != nil {
		log.Println("upload error", err)
		return c.JSON(fiber.Map{
			"status":  500,
			"message": "Server error",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  201,
		"message": "Deleted Successfully",
		"data":    nil,
	})
}
