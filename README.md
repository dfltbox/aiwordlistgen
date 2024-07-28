# aiwordlistgen
A password wordlist generator powered by ollama

## Example Usage
```go run . -m llama3.1 -q "includes the name jack" -a 5 -i 5 output.txt```

## Flags
-u The ollama API URL

-m The model you plan to use

-q The query for the model

-a The amount of passwords to generate in each request

-i The amount of times the password generating loops