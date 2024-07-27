package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/imroc/req/v3"
)

const url = "http://localhost:11434/api/generate"
const model = "llama3.1"
const query = "names with 'jack'"
const prompt = "Generate passwords that contain" + query + ". Return 10 passwords in a valid array. "
const iterations = 1

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

func clean(resp string) []string {
	cleanResp := strings.Trim(resp, "[]")
	cleanResp = strings.Trim(cleanResp, "\n")
	words := strings.Split(cleanResp, ",")

	var result []string
	for _, word := range words {
		cleanedElement := strings.TrimSpace(strings.Trim(word, "\""))
		result = append(result, cleanedElement)
	}
	return result
}

func appendtoOutput(text string) {
	outputFile := os.Args[1]
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	if _, err := file.WriteString(text); err != nil {
		log.Fatalf("failed to append text to file: %s", err)
	}
	defer file.Close()
}

func generate() {
	client := req.C()

	Prompt := &PromptT{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		System: "You are a password generating bot. Do not generate anything else. Make sure it is a valid array that can be used in a script. The array must follow the following format: ['password1','password2'].",
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
	//for debug
	response.Response = "['JackHammer88', 'JackSparrow123', 'JackDaniels77', 'JacksOrBetter11', 'JackOfAllTrades22', 'TheJackalope99', 'BlackJack777', 'CardSharkJack44', 'HighRollerJack55', 'LuckyJack999']"
	words := clean(response.Response)
	var appendedOutput string
	for _, word := range words {
		fmt.Println("Generated word", word)
		appendedOutput = appendedOutput + word + "\n"
	}
	appendedOutput = strings.ReplaceAll(appendedOutput, "'", "")
	appendtoOutput(appendedOutput)
	fmt.Println(response.Response)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a file to output stuff to!")
	}
	generate()
}
