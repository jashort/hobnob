package main

import (
	"github.com/alecthomas/kong"
	"hobnob/internal"
)

var CLI struct {
	Add struct {
		Name string   `arg:"" name:"name" help:"Name or alias" type:"string"`
		Note []string `arg:"" name:"note" help:"Note" type:"string"`
	} `cmd:"" help:"Add note"`
	About struct {
		Name string `arg:"" name:"name" help:"Name or alias" type:"string"`
	} `cmd:"" help:"List all notes about person"`
	Search struct {
		SearchString []string `arg:"" name:"string" help:"String to search for (case insensitive)" type:"string"`
		name         string   `short:"n" name:"name" help:"Only search person" type:"string" optional:"" placeholder:"Name|Alias"`
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
		ctx.Kong.Fatalf("Error loading data: %s", err)
	}
	switch ctx.Command() {
	case "search <search>":
		result = internal.CmdSearch(CLI.Search.SearchString, CLI.Search.name, &data)
	case "add <name> <note>":
		result = internal.CmdAdd(CLI.Add.Name, CLI.Add.Note, &data)
		err = internal.Save("data.json", data.Actions)
	case "undo":
		result = internal.CmdUndo(&data)
		err = internal.Save("data.json", data.Actions)
	case "history":
		result = internal.CmdHistory(&data)
	case "alias <name> <alias>":
		result, err = internal.CmdAlias(CLI.Alias.Name, CLI.Alias.Alias, &data)
		if err == nil {
			err = internal.Save("data.json", data.Actions)
		}
	case "aliases":
		result = internal.CmdAliases(&data)
	case "contacts":
		result = internal.CmdContacts(&data)
	case "about <name>":
		result = internal.CmdAbout(CLI.About.Name, &data)
	case "search <string>":
		result = internal.CmdSearch(CLI.Search.SearchString, CLI.Search.name, &data)
	default:
		ctx.Fatalf(ctx.Command())
	}
	if err != nil {
		ctx.Fatalf("%s", err)
	}
	println(result)
}
