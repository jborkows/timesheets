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
	ProjectRoot string
	Config      *model.Config
	Writer      io.Writer
}

type Controller struct {
	configs          *ControllerConfig
	didChangeReactor func(*messages.TextDocumentDidChangeNotification, *Controller)
	didSaveReactor   func(*messages.DidSaveTextDocumentNotification, *Controller)
	parser           *model.Parser
}

func NewController(c *ControllerConfig) *Controller {

	parser := model.Parser{
		HolidayClassifier: func(a *model.DateInfo) bool { return c.Config.IsHoliday(a) },
		IsCategory:        func(text string) bool { return c.Config.IsCategory(text) },
		IsTask:            func(text string) bool { return c.Config.IsTask(text) },
	}

	return &Controller{
		configs:          c,
		didChangeReactor: model.Debounce2(reactOnChange, time.Duration(1)*time.Second),
		didSaveReactor:   model.Debounce2(reactOnSave, time.Duration(1)*time.Second),
		parser:           &parser,
	}
}

func (self *Controller) onChange(msg *messages.TextDocumentDidChangeNotification) {
	self.didChangeReactor(msg, self)
}

func (self *Controller) onSave(msg *messages.DidSaveTextDocumentNotification) {
	self.didSaveReactor(msg, self)
}

type parsingTextParams struct {
	text string
	uri  string
}
type lineError struct {
	lineNumber int
	lineLength int
	err        error
}

func (self *Controller) parseText(input parsingTextParams) []model.WorkItem {

	var workItems []model.WorkItem = nil
	var errors []lineError = nil

	//TODO: pass date as parameter, extract from URI
	dateInfo := model.DateInfo{Value: "2021-01-05"}
	parseLine := self.parser.ParseLine(dateInfo)
	lines := strings.Split(input.text, "\n")
	for counter, line := range lines {
		if counter == len(lines)-1 && line == "" {
			continue
		}
		workItem, err := parseLine(line)
		if err != nil {
			errors = append(errors, lineError{lineNumber: counter, lineLength: len(line), err: err})
		} else {
			workItems = append(workItems, workItem)
		}
	}
	log.Printf("Parsed %+v items", workItems)
	log.Printf("Parsed %+v errors", errors)
	self.notifyAboutErrors(errors, input.uri)
	return workItems

}

func reactOnChange(msg *messages.TextDocumentDidChangeNotification, c *Controller) {
	text := msg.Params.ContentChanges[0].Text
	//TODO: use workItems
	//TODO: extract date from URI and pass it to parser
	c.parseText(parsingTextParams{text, msg.Params.TextDocument.URI})
	log.Println("Received didChange notification: ", msg.Params.TextDocument.URI)
}

func reactOnSave(msg *messages.DidSaveTextDocumentNotification, c *Controller) {
	log.Println("Received didSave notification: ", msg.Params.TextDocument.URI)
}

func (self *Controller) notifyAboutErrors(params []lineError, uri string) {
	var diagnostics []messages.Diagnostic = []messages.Diagnostic{}
	for _, param := range params {

		var errorMessage string
		switch param.err {
		case model.ErrEmptyLine:
			errorMessage = "Empty line"
		case model.ErrInvalidCategory:
			errorMessage = "Invalid category. Possible categories: " + strings.Join(self.configs.Config.PossibleCategories(), ", ")
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
				Start: messages.Position{Line: param.lineNumber, Character: 0},
				End:   messages.Position{Line: param.lineNumber, Character: 0},
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

	_, err := self.configs.Writer.Write([]byte(reply))
	return err

}
