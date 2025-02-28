package lspserver

import (
	"io"
	"log"
	"time"

	messages "github.com/jborkows/timesheets/internal/lspmessages"
	"github.com/jborkows/timesheets/internal/model"
	"github.com/jborkows/timesheets/internal/rpc"
)

type ControllerConfig struct {
	ProjectRoot string
	Config      *model.Config
}

type Controller struct {
	configs          *ControllerConfig
	didChangeReactor func(*messages.TextDocumentDidChangeNotification, *ControllerConfig)
	didSaveReactor   func(*messages.DidSaveTextDocumentNotification, *ControllerConfig)
}

func NewController(c *ControllerConfig) *Controller {
	return &Controller{
		configs:          c,
		didChangeReactor: model.Debounce2(reactOnChange, time.Duration(1)*time.Second),
		didSaveReactor:   model.Debounce2(reactOnSave, time.Duration(1)*time.Second),
	}
}

func (self *Controller) onChange(msg *messages.TextDocumentDidChangeNotification) {
	self.didChangeReactor(msg, self.configs)
}

func (self *Controller) onSave(msg *messages.DidSaveTextDocumentNotification) {
	self.didSaveReactor(msg, self.configs)
}

func reactOnChange(msg *messages.TextDocumentDidChangeNotification, c *ControllerConfig) {
	log.Println("Received didChange notification: ", msg.Params.TextDocument.URI)
}

func reactOnSave(msg *messages.DidSaveTextDocumentNotification, c *ControllerConfig) {
	log.Println("Received didSave notification: ", msg.Params.TextDocument.URI)
}

func writeResponse(writer io.Writer, msg any) error {
	reply := rpc.EncodeMessage(msg)

	_, err := writer.Write([]byte(reply))
	return err

}
