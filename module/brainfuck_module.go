// brainfuck_module.go
package module

import (
	"fmt"
	"github.com/krautchan/gbt/module/api"
	"github.com/krautchan/gbt/module/api/brainfuck"
	"github.com/krautchan/gbt/net/irc"
)

type BrainfuckModule struct {
	api.ModuleApi
}

func NewBrainfuckModule() *BrainfuckModule {
	return &BrainfuckModule{}
}

func (self *BrainfuckModule) Load() error {
	return nil
}

func (self *BrainfuckModule) GetCommands() map[string]string {
	return map[string]string{
		"bf": "SOURCE [INPUT] - Runs the given Brainfuck SOURCE with the given INPUT"}
}

func (self *BrainfuckModule) ExecuteCommand(cmd string, params []string, ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage) {
	if len(params) == 0 {
		return
	}

	source := params[0]
	input := ""

	if len(params) > 1 {
		input = params[1]
	}

	bf := brainfuck.NewBrainfuckInterpreter(source, input)
	output, err := bf.Start()
	if err != nil {
		c <- self.Reply(ircMsg, fmt.Sprintf("Error: %v", err.Error()))
		return
	}

	if len(output) > 0 {
		c <- self.Reply(ircMsg, string(output))
	}
}
