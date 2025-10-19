package aiagent

import "github.com/tousart/browser/models"

type AIAgent interface {
	DoPromptWithContent(message string) (*models.Answer, error)
}
