package main

import (
	"cms-octo-chat-api/model"
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {

	stmts, err := gormschema.New("postgres").Load(&model.User{}, &model.UserMatrix{}, &model.Conversation{}, &model.Message{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
