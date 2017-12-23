package labstack

import (
	"encoding/base64"
	"io/ioutil"
	"path/filepath"

	"github.com/dghubble/sling"
	"github.com/labstack/gommon/log"
)

type (
	// Jet defines the LabStack jet service.
	Jet struct {
		sling  *sling.Sling
		logger *log.Logger
	}

	// JetMessage defines the jet message.
	JetMessage struct {
		inlines     []string
		attachments []string
		Time        string     `json:"time,omitempty"`
		To          string     `json:"to,omitempty"`
		From        string     `json:"from,omitempty"`
		Subject     string     `json:"subject,omitempty"`
		Body        string     `json:"body,omitempty"`
		Inlines     []*jetFile `json:"inlines,omitempty"`
		Attachments []*jetFile `json:"attachments,omitempty"`
		Status      string     `json:"status,omitempty"`
	}

	jetFile struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Content string `json:"content"`
	}

	// JetError defines the jet error.
	JetError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

func NewJetMessage(to, from, subject string) *JetMessage {
	return &JetMessage{
		To:      to,
		From:    from,
		Subject: subject,
	}
}

func (m *JetMessage) addInlines() error {
	for _, inline := range m.inlines {
		data, err := ioutil.ReadFile(inline)
		if err != nil {
			return err
		}
		m.Inlines = append(m.Inlines, &jetFile{
			Name:    filepath.Base(inline),
			Content: base64.StdEncoding.EncodeToString(data),
		})
	}
	return nil
}

func (m *JetMessage) addAttachments() error {
	for _, attachment := range m.attachments {
		data, err := ioutil.ReadFile(attachment)
		if err != nil {
			return err
		}
		m.Inlines = append(m.Attachments, &jetFile{
			Name:    filepath.Base(attachment),
			Content: base64.StdEncoding.EncodeToString(data),
		})
	}
	return nil
}

func (m *JetMessage) AddInline(path string) {
	m.inlines = append(m.inlines, path)
}

func (m *JetMessage) AddAttachment(path string) {
	m.attachments = append(m.attachments, path)
}

// Send sends the jet message.
func (e *Jet) Send(m *JetMessage) (*JetMessage, error) {
	if err := m.addInlines(); err != nil {
		return nil, err
	}
	if err := m.addAttachments(); err != nil {
		return nil, err
	}
	em := new(JetMessage)
	ee := new(JetError)
	_, err := e.sling.Post("").BodyJSON(m).Receive(em, ee)
	if err != nil {
		return nil, err
	}
	if ee.Code == 0 {
		return em, nil
	}
	return nil, ee
}

func (e *JetError) Error() string {
	return e.Message
}
