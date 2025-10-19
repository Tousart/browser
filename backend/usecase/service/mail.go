package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

type MailService struct {
	webDriver selenium.WebDriver
}

func CreateMailService(webDriver selenium.WebDriver) *MailService {
	return &MailService{
		webDriver: webDriver,
	}
}

func (ms *MailService) Mail() ([]string, error) {
	ms.webDriver.Get("https://mail.yandex.ru")

	// Ожидаем, пока появятся письма
	err := ms.webDriver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		elems, _ := wd.FindElements(
			selenium.ByCSSSelector,
			"div.qa-MessagesListItemWrap",
		)
		return len(elems) > 0, nil
	}, 15*time.Second)
	if err != nil {
		return nil, fmt.Errorf("ошибка ожидания писем: %v", err)
	}

	// Находим письма
	// "[data-testid='messages-list_message-item_content']"
	rows, err := ms.webDriver.FindElements(
		selenium.ByCSSSelector,
		"div.qa-MessagesListItemWrap",
	)
	if err != nil || len(rows) == 0 {
		return nil, fmt.Errorf("писем не найдено")
	}

	emails := []string{}
	for i := range min(len(rows), 10) {
		subject, err := rows[i].FindElement(selenium.ByCSSSelector, "[data-testid='messages-list_subject']")
		if err != nil {
			log.Printf("не поймал субъекта %d: %v\n", i+1, err)
			continue
		}
		subjectText, _ := subject.Text()

		sender, err := rows[i].FindElement(selenium.ByCSSSelector, "[data-testid='message-common_sender-name']")
		if err != nil {
			log.Printf("не поймал отправителя %d: %v\n", i+1, err)
			continue
		}
		senderText, _ := sender.Text()

		emails = append(emails, fmt.Sprintf("От %s пришло письмо номер: %d с такой темой: %s", senderText, i+1, subjectText))
	}

	if len(emails) == 0 {
		return nil, errors.New("не добавлено ни одного письма")
	}

	return emails, nil
}
