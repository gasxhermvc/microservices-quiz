package domain

type EmailMessengerUseCase interface {
	SetTransaction(transId string)
	SendEmail(provider *EmailProvider) UseCaseResult
}

type EmailMessengerRepository interface{}

type EmailConfig struct {
	Server                  string `json:"server"`
	Port                    int    `json:"port"`
	EnableSSL               bool   `json:"enableSSL"`
	SendWithCredential      bool   `json:"sendWithCredential"`
	Skip                    bool   `json:"skip"`
	DefaultCredential       bool   `json:"defaultCredential"`
	Username                string `json:"username"`
	Password                string `json:"password"`
	SenderAddress           string `json:"sender_address"`
	FromParameter           string `json:"from_parameter"`
	ToParameter             string `json:"to_parameter"`
	CcParameter             string `json:"cc_parameter"`
	BccParameter            string `json:"bcc_parameter"`
	SubjectParameter        string `json:"subject_parameter"`
	BodyParameter           string `json:"body_parameter"`
	PriorityParameter       string `json:"priority_parameter"`
	AttachmentFileParameter string `json:"attachment_parameter"`
}

type EmailProvider struct {
	Config    EmailConfig
	Parameter EmailQueueParameter
}

type EmailQueueParameter struct {
	From       string                 `json:"from"`
	To         string                 `json:"to"`
	Cc         string                 `json:"cc"`
	Bcc        string                 `json:"bcc,omitempty"`
	Subject    string                 `json:"subject"`
	Body       string                 `json:"body"`
	Priority   string                 `json:"priority"`
	IsHtml     bool                   `json:"is_html"`
	Attachment []EmailQueueAttachment `json:"attachment,omitempty"`
}

type EmailQueueAttachment struct {
	FileId      string
	Filename    string
	MimeType    string
	Extension   string
	Size        float64
	DownloadUrl string
}
