package prompts

import (
	"fmt"
	"github/erastusk/tracer/internal/types"
	"github/erastusk/tracer/internal/utils"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func GetUserPrompt(z []string) string {
	var role string
	qs := []*survey.Question{
		{
			Name: "Debug Options",
			Prompt: &survey.Select{
				Message: "Select Debugging Option:",
				Options: z,
			},
		},
	}
	err := survey.Ask(qs, &role)
	if err != nil {
		// Handle cobra errors and Ctl - c commands
		close(utils.Done)
	}
	return role
}

func GetUserPromptSingle(prmpt string, p bool, a string) string {
	ans := ""

	if p {
		prompt := &survey.Password{
			Message: prmpt,
		}
		survey.AskOne(prompt, &ans)
	} else {
		prompt := &survey.Input{
			Message: prmpt,
			Default: a,
		}
		survey.AskOne(prompt, &ans)
	}
	return ans
}

func GetPrompts(a types.PromptOptions, options []string) (types.PromptOptions, error) {
	var result types.PromptOptions

	for _, r := range options {
		if strings.Contains(r, "password") {
			result.Password = GetUserPromptSingle(r, true, a.Password)
		}
		if strings.Contains(r, "endpoint") {
			result.Endpoint = GetUserPromptSingle(r, false, a.Endpoint)
		}
		if strings.Contains(r, "username") {
			result.Username = GetUserPromptSingle(r, false, a.Username)
		}
	}

	// Check for empty required fields
	if result.Endpoint == "" || result.Username == "" || result.Password == "" {
		return result, fmt.Errorf("required field/s empty")
	}
	return result, nil
}
