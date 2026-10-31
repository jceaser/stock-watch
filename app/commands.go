//***********************************************************************************************100
/*
Functions related to running commands as well as some standard system level commands
*/

package app

import (
	//standard modules
	"fmt"
	"maps"
	"slices"
	"strings"

	//personal modules
	"github.com/jceaser/stock-watch/common"

	//third party modules
	"github.com/kballard/go-shellquote"
)

//****************************************************************************80

var (
	HelpListFormat = "%15s  %s\n"
	HelpCmdFormat =  "%15s %-16s : %s\n"
)

type RunnableCommand func(any, Command, string,...string) string

type Command struct {
	Word string
	Func RunnableCommand
	Option string
	Description string
}

type Commands struct {
	List map[string]Command
}

func (self *Commands) A(word string, foo RunnableCommand, opt, desc string) {
	c := Command{strings.ToLower(word), foo, opt, desc}
	self.List[word] = c
}

func (self Commands) Run(state any, request string) string {
	requestAsWords, err := shellquote.Split(request)
	if err!=nil {
		return fmt.Sprintf("Parse Error: %v", err)
	}
	//requestAsWords := strings.Split(request, " ")\

	command := strings.ToLower(requestAsWords[0])

	justKeys := slices.Collect(maps.Keys(Cmd.List))
	slices.Sort(justKeys)
	for _, key := range justKeys {
		cmd := Cmd.List[key]
		if cmd.Word[0] >= 'a' && strings.HasPrefix(strings.ToLower(cmd.Word), command) {
			if cmd.Func == nil {
				continue
			}
			action := cmd.Func
			output := action(state, cmd, cmd.Option, requestAsWords[1:]...)
			return output
		}
	}
	return ""
}

// ***************************************************************************80

var Cmd Commands

// MARK - Stripting stuff

/* basic function to mark an action as not implemented yet */
func DoNoImpl (state any, command Command, flag string, options ...string) string {
	return "Not implimented"
}

/* utility action to echo out text */
func DoEcho(state any, command Command, flag string, options ...string) string {
	return fmt.Sprintf("\n%v\n", options)
}

/* overwrite this to have output go somewhere else */
func Out(format string, buffer *strings.Builder, a ...any) {
	fmt.Fprintf(buffer, format, a...)
}

func StandardHelpHeader(cmd Command) string {
	o := new(strings.Builder)
	Out("\n%s - %s\n\n", o, cmd.Word, cmd.Description)
	Out(HelpCmdFormat, o, "Command", "Flags", "Description")
	Out(HelpCmdFormat, o, "-------", "-----", "-----------")
	return o.String()
}

/* this is an action to print out the help document for all known commands */
func DoHelp(state any, command Command, flag string, options ...string) string {
	o := new(strings.Builder)
	if flag == "help" {
		cmd := Cmd.List["Help"]
		Out(StandardHelpHeader(cmd), o)
		Out(HelpCmdFormat, o, "help", "", "General Help, lists all commands")
		Out(HelpCmdFormat, o, "help", "<cmd>", "Help on a specific command")
		return o.String()
	}
	if len(options) == 0 {
		Out("\nHelp Manual\n\n", o)
		Out(HelpListFormat, o, "Command", "Description")
		Out(HelpListFormat, o, "----------", "----------")

		names := common.SortStringKeys[Command](Cmd.List)
		for _, n := range names {
			cmd := Cmd.List[n]
			Out(HelpListFormat, o, cmd.Word, cmd.Description)
		}
	} else {
		//for _, cmd := range Cmd.List {

		justKeys := slices.Collect(maps.Keys(Cmd.List))
		slices.Sort(justKeys)
		for _, key := range justKeys {
			cmd := Cmd.List[key]
			if strings.HasPrefix(strings.ToLower(cmd.Word), options[0]) {
				if cmd.Func == nil {
					Out(StandardHelpHeader(cmd), o)
					Out(fmt.Sprintf(HelpCmdFormat, cmd.Word, "", cmd.Description), o)
				} else {
					foo := cmd.Func
					Out(foo(state, cmd, "help", ""), o)
					break
				}
			}
		}
	}
	return o.String()
}
