package skill

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chromedp/chromedp"
)

// BrowserSkill provides browser control capabilities.
type BrowserSkill struct{}

func (s *BrowserSkill) FromJSON(jsonStr string) error {
	return nil
}

func (s *BrowserSkill) ToJSON() (string, error) {
	return "", nil
}

func (s *BrowserSkill) GetName() string { return "browser_manager" }

func (s *BrowserSkill) GetDescName() string {
	return "Browser Operations"
}

func (s *BrowserSkill) GetDescription() string {
	return "Browser operations: goto (visit), click (click), fill (input), observe (observe page structure/decision making), screenshot (take a screenshot)."
}

func (s *BrowserSkill) GetParameters() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action":   map[string]interface{}{"type": "string", "enum": []string{"goto", "click", "fill", "observe", "screenshot"}},
			"url":      map[string]interface{}{"type": "string"},
			"selector": map[string]interface{}{"type": "string", "description": "CSS selector"},
			"value":    map[string]interface{}{"type": "string", "description": "Input content"},
		},
		"required": []string{"action"},
	}
}

func (s *BrowserSkill) Execute(ctx context.Context, args string) (string, error) {
	var input struct {
		Action   string `json:"action"`
		URL      string `json:"url"`
		Selector string `json:"selector"`
		Value    string `json:"value"`
	}
	json.Unmarshal([]byte(args), &input)

	// This assumes that the context already contains the chromedp context.

	switch input.Action {
	case "goto":
		err := chromedp.Run(ctx, chromedp.Navigate(input.URL))
		if err != nil {
			return "navigation failed: " + err.Error(), err
		}
		return "successfully navigated to " + input.URL, nil

	case "observe":
		// Core function: extract key elements from the page for AI decision making
		var title string
		var description string
		// Simple logic: get the title and descriptions of all interactive elements
		err := chromedp.Run(ctx,
			chromedp.Title(&title),
			// The selector here can be refined as needed to only grab button, input, and a tags
			chromedp.Evaluate(`
				Array.from(document.querySelectorAll('button, input, a, [role="button"]'))
					.map(el => {
						return {
							tag: el.tagName,
							text: el.innerText || el.placeholder || el.value,
							id: el.id,
							name: el.getAttribute('name')
						}
					})
			`, &description),
		)
		if err != nil {
			return "observation failed", err
		}
		return fmt.Sprintf("Current page title: %s\nKey elements: %s", title, description), nil

	case "fill":
		err := chromedp.Run(ctx, chromedp.SendKeys(input.Selector, input.Value))
		if err != nil {
			return "fill failed", err
		}
		return "content entered in " + input.Selector, nil

	case "click":
		err := chromedp.Run(ctx, chromedp.Click(input.Selector))
		if err != nil {
			return "click failed", err
		}
		return "successfully clicked " + input.Selector, nil

	default:
		return "unknown action", nil
	}
}
