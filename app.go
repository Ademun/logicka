package main

import (
	"context"
	"fmt"
	"logicka/lib"
	"logicka/lib/visitor"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) CalculateTruthTable(expression string, variables map[string]bool) ([]visitor.TruthTableEntry, error) {
	if strings.TrimSpace(expression) == "" {
		return nil, fmt.Errorf("expression cannot be empty")
	}

	truthTableData, err := lib.GenerateTruthTable(expression, variables)
	if err != nil {
		return nil, fmt.Errorf("error generating truth table: %v", err)
	}

	return truthTableData, nil
}

func (a *App) ExtractVariables(expression string) ([]string, error) {
	vars := lib.ExtractVariables(expression)
	return vars, nil
}
