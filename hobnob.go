package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"hobnob/internal"
	"log"
)

var CLI struct {
	Add struct {
		Name string   `arg:"" name:"name" help:"Name or alias" type:"string"`
		Note []string `arg:"" name:"note" help:"Note" type:"string"`
	} `cmd:"" help:"Add note"`
	From struct {
		Name string `arg:"" name:"name" help:"Name or alias" type:"string"`
	} `cmd:"" help:"List all notes from person"`
	Search struct {
		SearchString []string `arg:"" name:"search string" help:"String to search for (case insensitive)" type:"string"`
	} `cmd:"" help:"Search notes"`
	Undo  struct{} `cmd:"" help:"Undo last action"`
	Alias struct {
		Name  string `arg:"" name:"name" help:"(Full) Name" type:"string"`
		Alias string `arg:"" name:"alias" help:"Alias" type:"string"`
	} `cmd:"" help:"Add alias"`
	Aliases  struct{} `cmd:"" help:"List all aliases"`
	Contacts struct{} `cmd:"" help:"List all contacts"`
	History  struct{} `cmd:"" help:"Show history"`
}

func main() {
	ctx := kong.Parse(&CLI)
	data, err := internal.LoadAll("data.json")
	result := ""
	if err != nil {
		panic(err)
	}
	switch ctx.Command() {
	case "search <search>":
		result = internal.CmdSearch(ctx, &data)
	case "add <name> <note>":
		result = internal.CmdAdd(ctx, &data)
		err = internal.Save("data.json", data.Actions)
	case "undo":
		result = internal.CmdUndo(ctx, &data)
		err = internal.Save("data.json", data.Actions)
	case "history":
		result = internal.CmdHistory(ctx, &data)
	case "alias <name> <alias>":
		result, err = internal.CmdAlias(ctx, &data)
		if err == nil {
			err = internal.Save("data.json", data.Actions)
		}
	case "aliases":
		result = internal.CmdAliases(ctx, &data)
	case "contacts":
		result = internal.CmdContacts(ctx, &data)
	case "from <name>":
		result = internal.CmdFrom(ctx, &data)
	default:
		fmt.Println("Unknown command")
		log.Fatal(ctx.Command())
	}
	if err != nil {
		panic(err)
	}
	println(result)
}
