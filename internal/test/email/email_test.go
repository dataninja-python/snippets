package email_test

type EmailAddress struct {
	Local     string
	Separator string
	Domain    string
}

// emailSeparator or the "at sign" separates the local from the domain portion of the email address.
const emailSeparator string = "@"
