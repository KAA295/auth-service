package services

import "fmt"

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (e *EmailService) Send(msg string) {
	fmt.Printf(`message "%v" sent\n`, msg)
}
