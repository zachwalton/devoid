package main

import (
	"context"
	"os"

	"github.com/zachwalton/devoid/cmd"
	_ "github.com/zachwalton/devoid/pkg/tui"
)

func main() {
	(cmd.Cmd).Run(context.Background(), os.Args)
}
