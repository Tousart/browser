package usecase

type MailService interface {
	Mail() ([]string, error)
}
