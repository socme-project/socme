package notifyme

import (
	ssmtp "net/smtp"
)

type smtp struct {
	from     string
	password string
	to       string
	url      string // smtp.gmail.com
	port     string // 587
}

func Smtp(from, password, to, url, port string) *smtp {
	// parse the url
	return &smtp{
		from:     from,
		password: password,
		to:       to,
		url:      url,
		port:     port,
	}
}

func (s *smtp) Notify(title, message string) error {
	msg := "From: " + s.from + "\n" +
		"To: " + s.to + "\n" +
		"Subject: " + title + "\n\n" +
		message

	err := ssmtp.SendMail(s.url+":"+s.port,
		ssmtp.PlainAuth("", s.from, s.password, s.url),
		s.from, []string{s.to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
