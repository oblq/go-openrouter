package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/revrost/go-openrouter"
)

func main() {
	client := openrouter.NewClient(os.Getenv("OPENROUTER_API_KEY"))

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Model: openrouter.GeminiFlashImage,
			Messages: []openrouter.ChatCompletionMessage{
				{
					Role:    openrouter.ChatMessageRoleUser,
					Content: openrouter.Content{Text: "Create a picture of a nano banana dish in a fancy restaurant with a Gemini theme"},
				},
			},
			Modalities: []openrouter.ChatCompletionModality{
				openrouter.ModalityImage,
			},
			ImageConfig: &openrouter.ImageConfig{
				AspectRatio: openrouter.AspectRatioPortrait_9_16,
			},
		},
	)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(resp.Choices) == 0 {
		log.Fatalln("no choices...")
	}

	message := resp.Choices[0].Message
	if len(message.Images) == 0 {
		log.Fatalln("image not generated")
	}

	imageBase64URL := message.Images[0].ImageURL.URL
	if len(imageBase64URL) == 0 {
		log.Fatalln("corrupted image")
	}

	if strings.HasPrefix(imageBase64URL, "data:image/") {
		parts := strings.Split(imageBase64URL, ",")
		if len(parts) > 1 {
			imageBase64URL = parts[1]
		}
	}

	var imagePngBytes []byte
	imagePngBytes, err = base64.StdEncoding.DecodeString(imageBase64URL)
	if err != nil {
		log.Fatalf("error decoding base64 image URL: %s", err)
	}

	filename := "generated_image.png"
	err = os.WriteFile(filename, imagePngBytes, 0644)
	if err != nil {
		log.Printf("Error saving image: %v", err)
		return
	}
	log.Printf("Image saved to: %s", filename)
}
