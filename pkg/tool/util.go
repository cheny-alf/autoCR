package tool

import "fmt"

func buildParam(prompt, content string) string {
	// add common prompt, make sure that AI returns in accordance with the specified format, like [25] xxxx
	defaultPrompt := "%s. You must return it in this format, like [25] if err ! = nil { . instead of [Line 25] if err ! = nil \n %s"
	return fmt.Sprintf(prompt, defaultPrompt, content)
}
