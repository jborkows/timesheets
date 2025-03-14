package lspserver

import (
	"fmt"
	"io"
	"log"
	"net/url"
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
	content          *content
}

func NewController(c *ControllerConfig) *Controller {

	return &Controller{
		service:          c.Service,
		writer:           c.Writer,
		didChangeReactor: model.Debounce2(reactOnChange, time.Duration(1)*time.Second),
		didSaveReactor:   model.Debounce2(reactOnSave, time.Duration(1)*time.Second),
		content:          newContent(),
	}
}

func (self *Controller) onChange(msg *messages.TextDocumentDidChangeNotification) {
	text := msg.Params.ContentChanges[0].Text
	self.content.put(msg.Params.TextDocument.URI, strings.Split(text, "\n"))
	self.didChangeReactor(msg, self)
}

func (self *Controller) onSave(msg *messages.DidSaveTextDocumentNotification) {
	self.changeContent(msg.Params.TextDocument)
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

func (self *Controller) onOpen(msg *messages.DidOpenTextDocumentNotification) {
	self.changeContent(msg.Params.TextDocument)
}

func (self *Controller) changeContent(textDocument messages.TextDocumentItem) {
	self.content.put(textDocument.URI, strings.Split(textDocument.Text, "\n"))
}

func (self *Controller) notifyAboutErrors(params []model.LineError, uri string) {
	var diagnostics []messages.Diagnostic = []messages.Diagnostic{}
	for _, param := range params {

		var errorMessage string
		switch param.Err {
		case model.ErrEmptyLine:
			continue
		case model.ErrInvalidCategory:
			errorMessage = "Invalid category. Possible categories: " + strings.Join(self.service.PossibleCategories(), ", ")
		case model.ErrInvalidTime:
			errorMessage = "Invalid time format. Use X.Y or XhYm (e.g., 1.5 or 1h30m)"
		default:
			errorMessage = param.Err.Error()
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

func (self *Controller) completion(request *messages.CompletionRequest) error {
	params := request.Params
	completions := self.completions(params.TextDocument.URI, params.Position)

	msg := messages.CompletionResponse{
		Response: response(request.Request),
		Result:   completions,
	}
	return self.writeResponse(msg)
}

func (self *Controller) completions(uri string, position messages.Position) []messages.CompletionItem {
	if position.Character > 5 {
		return []messages.CompletionItem{}
	}
	completions := []messages.CompletionItem{}
	for _, category := range self.service.PossibleCategories() {
		completions = append(completions, messages.CompletionItem{
			Label:  category,
			Detail: "Category",
		})
	}
	return completions
}

func (self *Controller) Hover(request *messages.HoverRequest) error {

	date, err := self.service.ParseDateFromName(request.Params.TextDocument.URI)
	if err != nil {
		return fmt.Errorf("Error getting date from file: %s for %s", err, request.Params.TextDocument.URI)
	}
	statistics, err := self.service.DailyStatistics(date)
	if err != nil {
		return fmt.Errorf("Error getting daily statistics: %s", err)
	}
	var builder strings.Builder
	builder.WriteString("Daily statistics for ")
	builder.WriteString(date.Format("2006-01-02"))
	builder.WriteString("\n")
	for _, stat := range statistics {
		source := stat.Dirty
		builder.WriteString(source.Category)
		builder.WriteString(": ")
		builder.WriteString(fmt.Sprintf("%02d:%02d", source.Hours, source.Minutes))
		builder.WriteString("\n")
	}
	builder.WriteString("End \n")

	msg := messages.HoverResponse{
		Response: response(request.Request),
		Result: messages.HoverResult{
			Contents: builder.String(),
		},
	}

	return self.writeResponse(msg)
}

func (self *Controller) Definition(request *messages.DefinitionRequest) error {
	date, err := self.service.ParseDateFromName(request.Params.TextDocument.URI)
	if err != nil {
		return fmt.Errorf("Error getting date from file: %s for %s", err, request.Params.TextDocument.URI)
	}
	output, err := self.service.ShowDailyStatistics(date)
	if err != nil {
		return fmt.Errorf("Error getting statistics for file: %s for %s", err, request.Params.TextDocument.URI)
	}

	uri := url.URL{Scheme: "file", Path: string(output)}
	msg := messages.DefinitionResponse{
		Response: response(request.Request),
		Result: &messages.Location{
			URI: uri.String(),
			Range: messages.Range{
				Start: messages.Position{Line: 0, Character: 0},
				End:   messages.Position{Line: 0, Character: 0},
			},
		},
	}
	return self.writeResponse(msg)

}

func (self *Controller) writeResponse(msg any) error {
	reply := rpc.EncodeMessage(msg)

	_, err := self.writer.Write([]byte(reply))
	return err

}
