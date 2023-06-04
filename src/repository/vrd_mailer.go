package repository

type Vrd_mailer struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	MailTo  string `json:"mailto"`
}

type Vrd_mailer_return struct {
	Message string `json:"message"`
}