package parser

import (
	"errors"
	"fmt"

	"github.com/google/shlex"
)

const (
	PREFIX    = "mqtt"
	SUBSCRIBE = "subscribe"
	PUBLISH   = "publish"
	SEND      = "send"
)

type Command struct {
	Command string
	Topic   string
	Value   string
}

func (c Command) String() (string, error) {
	if c.Command != SEND {
		return "", fmt.Errorf("can not send command %v", c.Command)
	}
	return fmt.Sprintf("%v %v %v", PREFIX, c.Topic, c.Value), nil
}

// ParseLine parses a line and returns the corresponding command.
func ParseLine(line string) (*Command, error) {
	// Use shell lexer for token splitting
	token, err := shlex.Split(line)
	if err != nil {
		return nil, fmt.Errorf("lexer error while parsing line, %v", err)
	}

	length := len(token)

	// No mqtt message. This is no error but there is no command either
	if length < 1 || token[0] != PREFIX {
		return nil, nil
	}

	if length < 2 {
		return nil, errors.New("no command found")
	}

	switch token[1] {
	case SUBSCRIBE:
		if length < 3 {
			return nil, errors.New("not enough parameter for subscribe command")
		}
		return &Command{
			Command: SUBSCRIBE,
			Topic:   token[2],
		}, nil
	case PUBLISH:
		if length < 4 {
			return nil, errors.New("not enough parameter for publish command")
		}
		return &Command{
			Command: PUBLISH,
			Topic:   token[2],
			Value:   token[3],
		}, nil
	default:
		return nil, fmt.Errorf("unknown command %v", token[1])
	}
}
