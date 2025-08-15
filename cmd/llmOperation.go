package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"google.golang.org/genai"
)

func ProcessWithLLM(clipboardData string) string {

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	prompt := fmt.Sprintf(`You are a smart text processor. I will provide a block of text. Your task is to output a JSON object with the following fields:

{
  "title": "A concise title summarizing the text in one line",
  "text": "The same text minimally formatted: remove extra spaces, unnecessary characters, trim leading/trailing spaces",
  "tag": "A single label describing the type of text, e.g., code, SQL query, plain text, credentials, notes, log, configuration, command, etc. Choose the most appropriate label."
}
⚠ Important:
- Do NOT include markdown formatting, code fences, or any extra text.  
- Output must be plain JSON only so it can be parsed programmatically.
Do not include anything other than the JSON object. Format it exactly like JSON so it can be parsed programmatically.

Here is the text to process:

"%s"`, clipboardData)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	return result.Text()
}

func CleanLLMOutput(llmOutput string) string {
	// Trim leading/trailing whitespace
	cleaned := strings.TrimSpace(llmOutput)

	// Remove ```json or ``` fences (any language or no language)
	re := regexp.MustCompile("(?s)```[a-zA-Z]*\\n?(.*?)```")
	matches := re.FindStringSubmatch(cleaned)
	if len(matches) > 1 {
		// Extract inner content
		cleaned = matches[1]
	}

	return strings.TrimSpace(cleaned)
}

func ParseLLMoutput(CleanllmOutput string) (*JSONformatClipboardData, error) {

	var result JSONformatClipboardData

	err := json.Unmarshal([]byte(CleanllmOutput), &result)

	if err != nil {
		return nil, err
	}

	return &result, nil

}
