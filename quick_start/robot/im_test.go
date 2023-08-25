package main

import (
	"fmt"
	"testing"
)

func TestListChatHistory(t *testing.T) {
	err := ListChatHistory("xxxx")
	if err != nil {
		fmt.Println(err)
	}
}
