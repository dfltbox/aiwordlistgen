package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/imroc/req/v3"
)

var url string
var model string
var query string
var count int
var prompt = "Generate passwords that contain" + query + ". Return" + string(count) + "passwords in a valid array. "
var iterations int
var completed int

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

func appendtoOutput(text string, outputFile string) {
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	if _, err := file.WriteString(text); err != nil {
		log.Fatalf("failed to append text to file: %s", err)
	}
	defer file.Close()
}

func generate() string {
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
	//response.Response = "['JackHammer88', 'JackSparrow123', 'JackDaniels77', 'JacksOrBetter11', 'JackOfAllTrades22']"
	words := clean(response.Response)
	var appendedOutput string
	for _, word := range words {
		appendedOutput = appendedOutput + word + "\n"
	}
	appendedOutput = strings.ReplaceAll(appendedOutput, "'", "")

	//fmt.Println(response.Response)

	return appendedOutput
}

func main() {
	completed = 0
	//flags and command args

	flag.StringVar(&url, "u", "http://localhost:11434/api/generate", "The ollama API URL")
	flag.StringVar(&model, "m", "", "The model you plan to use")
	flag.StringVar(&query, "q", "", "The query for the model")
	flag.IntVar(&iterations, "a", 5, "The amount of passwords to generate in each request")
	flag.IntVar(&count, "i", 5, "The amount of times the password generating loops")

	flag.Parse()

	if len(os.Args) < 2 {
		log.Fatal("Please provide a file to output stuff to!")
	}
	if count < 1 {
		log.Fatal("You have to have at least 1 password each iteration!\nExample: -c 5")
	}
	if iterations < 1 {
		log.Fatal("You have to iterate at least once!\nExample: -i 5")
	}
	if query == "" {
		log.Fatal("Please provide a query\nExample: -q 'includes the name jack'")
	}
	if model == "" {
		log.Fatal("Please provide a model\nExample: -m llama3.1")
	}
	for i := 0; i < iterations; i++ {
		appenedOutput := generate()
		oFile := flag.Args()[0]
		appendtoOutput(appenedOutput, oFile)
		completed += count
		fmt.Println("Generated", completed, "/", iterations*count)
	}
	fmt.Println("Done!")

}
