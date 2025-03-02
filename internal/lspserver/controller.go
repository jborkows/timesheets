package lspserver

import (
	"io"
	"log"
	"strings"
	"time"

	messages "github.com/jborkows/timesheets/internal/lspmessages"
	"github.com/jborkows/timesheets/internal/model"
	"github.com/jborkows/timesheets/internal/rpc"
)

type ControllerConfig struct {
	Service *model.Service
	Writer  io.Writer
}

type Controller struct {
	didChangeReactor func(*messages.TextDocumentDidChangeNotification, *Controller)
	didSaveReactor   func(*messages.DidSaveTextDocumentNotification, *Controller)
	service          *model.Service
	writer           io.Writer
}

func NewController(c *ControllerConfig) *Controller {

	return &Controller{
		service:          c.Service,
		writer:           c.Writer,
		didChangeReactor: model.Debounce2(reactOnChange, time.Duration(1)*time.Second),
		didSaveReactor:   model.Debounce2(reactOnSave, time.Duration(1)*time.Second),
	}
}

func (self *Controller) onChange(msg *messages.TextDocumentDidChangeNotification) {
	self.didChangeReactor(msg, self)
}

func (self *Controller) onSave(msg *messages.DidSaveTextDocumentNotification) {
	self.didSaveReactor(msg, self)
}

func reactOnChange(msg *messages.TextDocumentDidChangeNotification, c *Controller) {
	text := msg.Params.ContentChanges[0].Text
	date, err := c.service.ParseDateFromName(msg.Params.TextDocument.URI)
	if err != nil {
		log.Printf("Error getting date from file: %s for %s", err, msg.Params.TextDocument.URI)
	}
	_, errors := c.service.ProcessForDraft(text, date)
	c.notifyAboutErrors(errors, msg.Params.TextDocument.URI)
	log.Println("Received didChange notification: ", msg.Params.TextDocument.URI, "representing", date)
}

func reactOnSave(msg *messages.DidSaveTextDocumentNotification, c *Controller) {
	text := msg.Params.Text
	date, err := c.service.ParseDateFromName(msg.Params.TextDocument.URI)
	if err != nil {
		log.Printf("Error getting date from file: %s for %s", err, msg.Params.TextDocument.URI)
	}
	_, errors := c.service.ProcessForSave(text, date)
	c.notifyAboutErrors(errors, msg.Params.TextDocument.URI)
	log.Println("Received didSave notification: ", msg.Params.TextDocument.URI)
}

func (self *Controller) notifyAboutErrors(params []model.LineError, uri string) {
	var diagnostics []messages.Diagnostic = []messages.Diagnostic{}
	for _, param := range params {

		var errorMessage string
		switch param.Err {
		case model.ErrEmptyLine:
			errorMessage = "Empty line"
		case model.ErrInvalidCategory:
			errorMessage = "Invalid category. Possible categories: " + strings.Join(self.service.PossibleCategories(), ", ")
		case model.ErrInvalidTime:
			errorMessage = "Invalid time format. Use X.Y or XhYm (e.g., 1.5 or 1h30m)"
		default:
			errorMessage = "Unknown error"
			log.Printf("Unknown error: %v", param)
		}
		diagnostic := messages.Diagnostic{
			Message:  errorMessage,
			Severity: 1,
			Range: messages.Range{
				Start: messages.Position{Line: param.LineNumber, Character: 0},
				End:   messages.Position{Line: param.LineNumber, Character: 0},
			},
		}
		diagnostics = append(diagnostics, diagnostic)
	}

	message := messages.PublishDiagnosticsNotification{
		Notification: messages.Notification{
			RPC:    "2.0",
			Method: "textDocument/publishDiagnostics",
		},
		Params: messages.PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diagnostics,
		},
	}
	log.Printf("Sending diagnostics for %s with %+v", uri, message)
	err := self.writeResponse(message)
	if err != nil {
		log.Printf("Error writing response for notifying about errors: %s", err)
	}
}

func (self *Controller) writeResponse(msg any) error {
	reply := rpc.EncodeMessage(msg)

	_, err := self.writer.Write([]byte(reply))
	return err

}
