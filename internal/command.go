package internal

import (
	"cmp"
	"fmt"
	"github.com/alecthomas/kong"
	"maps"
	"slices"
	"strings"
	"time"
)

func CmdAdd(ctx *kong.Context, data *Data) string {
	record := Note{
		Name:      data.LookupName(ctx.Args[1]),
		Note:      strings.Join(ctx.Args[2:], " "),
		Timestamp: time.Now().Local(),
	}

	AddNote(data, record)
	return record.String()
}

func CmdSearch(ctx *kong.Context, data *Data) string {
	searchString := strings.ToLower(strings.Join(ctx.Args[1:], " "))
	output := strings.Builder{}
	for _, note := range data.Notes {
		if strings.Contains(strings.ToLower(note.Note), searchString) {
			output.WriteString(fmt.Sprintf("%s\n\n", note))
		}
	}
	return output.String()
}

func CmdFrom(ctx *kong.Context, data *Data) string {
	searchString := strings.ToLower(strings.Join(ctx.Args[1:], " "))

	name := strings.ToLower(data.LookupName(searchString))
	output := strings.Builder{}
	for _, note := range data.Notes {
		if strings.ToLower(note.Name) == name {
			output.WriteString(fmt.Sprintf("%s\n\n", note))
		}
	}
	return output.String()
}

func CmdAlias(ctx *kong.Context, data *Data) (string, error) {
	record := Alias{
		Name:  ctx.Args[1],
		Alias: strings.Join(ctx.Args[2:], " "),
	}
	err := AddAlias(data, record)
	return fmt.Sprintf("Added alias from name %s to alias %s", record.Name, record.Alias), err
}

func CmdAliases(ctx *kong.Context, data *Data) string {
	aliasCmp := func(a, b Alias) int {
		return cmp.Compare(strings.ToLower(a.Alias), strings.ToLower(b.Alias))
	}
	slices.SortFunc(data.Aliases, aliasCmp)
	output := strings.Builder{}
	output.WriteString("Aliases:\n")
	for _, element := range data.Aliases {
		output.WriteString(fmt.Sprintf("  Alias %s -> Name %s\n", element.Alias, element.Name))
	}
	return output.String()
}

func CmdContacts(ctx *kong.Context, data *Data) string {
	people := make(map[string]*Person)
	output := strings.Builder{}
	for _, alias := range data.Aliases {
		lowerName := strings.ToLower(alias.Name)
		if _, exists := people[lowerName]; !exists {
			people[lowerName] = &Person{Name: alias.Name}
		}
		people[lowerName].Aliases = append(people[lowerName].Aliases, alias.Alias)
	}

	for _, element := range data.Notes {
		lowerName := strings.ToLower(element.Name)
		if _, exists := people[lowerName]; !exists {
			people[lowerName] = &Person{Name: element.Name}
		}
	}
	personCmp := func(a, b *Person) int {
		return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	}
	pList := slices.Collect(maps.Values(people))
	slices.SortFunc(pList, personCmp)
	for _, person := range pList {
		output.WriteString(fmt.Sprintf("  %s\n", person))
	}
	return output.String()
}

// CmdUndo Because this is a CLI utility, "undo" removes the last thing from the list of
// actions, but does not worry about cleaning up the rest of the data structure because
// the data file will be reloaded next time the command is run.
func CmdUndo(ctx *kong.Context, data *Data) string {
	if len(data.Actions) > 0 {
		data.Actions = data.Actions[:len(data.Actions)-1]
		return fmt.Sprintf("Removed 1 record, %d remaining\n", len(data.Actions))
	} else {
		return "No records to remove"
	}
}

func CmdHistory(ctx *kong.Context, data *Data) string {
	output := strings.Builder{}
	for _, action := range data.Actions {
		output.WriteString(fmt.Sprintf("%s %s\n  %s\n", action.Timestamp, action.Action, action.Data))
	}
	return output.String()
}
