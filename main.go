package main

import (
	"encoding/json"
	"fmt"
	"io"
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

	type ResponseT struct {
		EvalDuration int    `json:"eval_duration"`
		EvalCount    int    `json:"eval_count"`
		Response     string `json:"response"`
	}

	Prompt := &PromptT{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		System: "You are a password generating bot. Do not generate anything else. Make sure it is a valid array that can be used in a script. Do not include newlines in your response. The array must follow the following format: ['password1','password2']",
	}

	resp, err := client.R().
		SetBody(Prompt).Post(url)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var response ResponseT

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	fmt.Println(response.Response)
}
