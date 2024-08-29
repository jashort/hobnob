package internal

import (
	"strings"
	"testing"
)

func TestCommands(t *testing.T) {
	data := Data{
		Aliases: nil,
		Notes:   nil,
		Actions: nil,
		Version: 1,
	}

	// Add some aliases
	res, _ := CmdAlias("Alice Smith", "alice", &data)
	if res != "Added alias from name Alice Smith to alias alice" {
		t.Error("Incorrect output from adding alias alice")
	}
	res, _ = CmdAlias("Bob Jones", "Bob", &data)
	if res != "Added alias from name Bob Jones to alias Bob" {
		t.Error("Incorrect output from adding alias Bob")
	}
	// Adding the same alias twice should return an error
	_, tmpErr := CmdAlias("Bob Jones", "Bob", &data)
	if tmpErr == nil {
		t.Error("Adding duplicate alias did not return error")
	}

	// Add some notes
	res = CmdAdd("Alice Smith", []string{"Alice note 1"}, &data)
	if !strings.Contains(res, "Alice Smith") {
		t.Error("Expected 'Alice Smith', got ", res)
	}

	// Using an alias should store the person's full name
	res = CmdAdd("Alice", []string{"Alice note 2"}, &data)
	if !strings.Contains(res, "Alice Smith") {
		t.Error("Expected 'Alice Smith', got ", res)
	}

	res = CmdAdd("bob", []string{"Bob", "note 1"}, &data)
	if !strings.Contains(res, "Bob note 1") {
		t.Error("Expected 'Bob note 1', got ", res)
	}
	if len(data.Actions) != 5 {
		t.Error("Expected 5 actions, got ", len(data.Actions))
	}
	if len(data.Aliases) != 2 {
		t.Error("Expected 2 aliases, got ", len(data.Aliases))
	}
	if len(data.Notes) != 3 {
		t.Error("Expected 3 notes, got ", len(data.Notes))
	}

	// Remove the last note
	res = CmdUndo(&data)
	if len(data.Actions) != 4 {
		t.Error("Expected 4 actions, got ", len(data.Actions))
	}

	res = CmdContacts(&data)
	if !strings.Contains(res, "Alice Smith  (alice)\n") || !strings.Contains(res, "Bob Jones  (Bob)\n") {
		t.Error("Incorrect output from contacts commands")
	}

	res = CmdHistory(&data)
	if !strings.Contains(res, "Alice Smith -> alice") ||
		!strings.Contains(res, `Bob Jones -> Bob`) ||
		!strings.Contains(res, "Alice note 1") {
		t.Error("Incorrect output from history command")
	}

	res = CmdAbout("alice", &data)
	if !strings.Contains(res, "Alice Smith") ||
		!strings.Contains(res, "Alice note 1") ||
		!strings.Contains(res, "Alice note 2") {
		t.Error("Incorrect output from from command")
	}

	res = CmdSearch([]string{"bob", "note", "1"}, "bob", &data)
	if !strings.Contains(res, "Bob note 1") {
		t.Error("Incorrect output from search command, got", res)
	}

	res = CmdAliases(&data)
	if !strings.Contains(res, "Aliases:") ||
		!strings.Contains(res, "Alias alice -> Name Alice Smith") ||
		!strings.Contains(res, "Alias Bob -> Name Bob Jones") {
		t.Error("Incorrect output from aliases command:", res)
	}
}
