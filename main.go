package main

import (
	"fmt"
	"log"

	"github.com/imroc/req/v3"
)

func main() {
	const url = "http://localhost:11434/api/generate"
	const model = "llama3.1"
	const query = "names with 'jack'"
	const prompt = "Generate passwords that contain" + query + ". Return 10 passwords in a valid array."

	client := req.C()

	type PromptT struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Format string `json:"format"`
		Stream bool   `json:"stream"`
		System string `json:"system"`
		Raw    bool   `json:"raw"`
	}

	Prompt := &PromptT{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		System: "You are a password generating bot. Do not generate anything else. Make sure it is a valid array that can be used in a script. Do not include newlines in your response.",
	}

	resp, err := client.R().
		SetBody(Prompt).Post(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
