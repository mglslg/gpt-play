package ds

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Stack struct {
	elements []*discordgo.Message
}

func NewStack() *Stack {
	return &Stack{elements: make([]*discordgo.Message, 0)}
}

func (s *Stack) Push(element *discordgo.Message) {
	s.elements = append(s.elements, element)
}

func (s *Stack) Pop() (*discordgo.Message, error) {
	if len(s.elements) == 0 {
		return nil, fmt.Errorf("stack is empty")
	}

	topIndex := len(s.elements) - 1
	topElement := s.elements[topIndex]
	s.elements = s.elements[:topIndex]

	return topElement, nil
}

func (s *Stack) Size() int {
	return len(s.elements)
}
