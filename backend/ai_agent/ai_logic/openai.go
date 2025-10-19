package ailogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/tousart/browser/models"
	"github.com/tousart/browser/usecase"
)

type AIAgent struct {
	openAIClient *openai.Client
	service      usecase.MailService
	roles        []string
}

func CreateAIAgent(apiKey string, service usecase.MailService) *AIAgent {
	client := openai.NewClient(apiKey)
	roles := []string{"'письма и почта'", "'работа и вакансии'", "'заказ и еда'"}
	return &AIAgent{
		openAIClient: client,
		service:      service,
		roles:        roles,
	}
}

func (ai *AIAgent) DoPromptWithContent(message string) (*models.Answer, error) {
	prompt := fmt.Sprintf(`Исходя из запроса, выбери только ОДНУ сферу деятельности в одинарных кавычках, которая соответствует запросу больше всего: %s, если не подходит ничего из перечисленного, то выбери 'не умею'.`,
		strings.Join(ai.roles, ", "))

	resp, err := ai.openAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: prompt},
				{Role: openai.ChatMessageRoleUser, Content: message},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("ai_agent ai_logic DoPromptWithContent error: %v", err)
	}

	content := resp.Choices[0].Message.Content

	var answer models.Answer
	switch content {
	case ai.roles[0]:
		err = ai.emailsToSpam(&answer)
	default:
		answer.AIAnswer = "smth"
	}

	if err != nil {
		return nil, fmt.Errorf("ai_agent ai_logic DoPromptWithContent error: %v", err)
	}

	return &answer, nil
}

func (ai *AIAgent) emailsToSpam(answer *models.Answer) error {
	emails, err := ai.service.Mail()
	if err != nil {
		return err
	}

	prompt := fmt.Sprintf(`Дан список писем, определи, какие из этих писем - спам (реклама, акции).
	Ответь списком номеров спам-писем через запятую (например: 1,3,4) и не пиши к нему никакого текста (если все письма нормальные пришли 0).
	Письма:
	%s`, strings.Join(emails, "\n"))

	resp, err := ai.openAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "Ты помощник, определяющий спам-письма."},
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
		},
	)

	if err != nil {
		return err
	}

	sliceOfEmailsNumbers := strings.Split(resp.Choices[0].Message.Content, ",")
	if len(sliceOfEmailsNumbers) == 1 && sliceOfEmailsNumbers[0] == "0" {
		answer.AIAnswer = "Последние 10 писем не похожи на спам"
	} else {
		// TODO отправление писем в спам по номерам
		answer.AIAnswer = strings.Join(sliceOfEmailsNumbers, ",")
	}

	return nil
}
