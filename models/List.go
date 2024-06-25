package models

import (
	"fmt"
	"strings"

	"github.com/andrwui/hyprsources/cli"
	"github.com/andrwui/hyprsources/cli/style"
	"github.com/andrwui/hyprsources/constants"
	"github.com/andrwui/hyprsources/constants/colors"
	"github.com/andrwui/hyprsources/helper"
)

type List struct {
	Prompt      string
	CursorPos   int
	SourceFiles []*SourceFile
}

func NewList(prompt string) *List {
	return &List{
		Prompt:      prompt,
		SourceFiles: make([]*SourceFile, 0),
	}
}

func (l *List) AddSourceToList(sourceFile SourceFile) *List {
	l.SourceFiles = append(l.SourceFiles, &sourceFile)
	return l
}

func (l *List) SetSources(sfs []*SourceFile) *List {
	l.SourceFiles = sfs
	return l
}

func (l *List) calculateMaxNameAndPathLen() (int, int) {
	var maxNameLen, maxPathLen int
	for _, sourceFile := range l.SourceFiles {
		if len(sourceFile.Name) > maxNameLen {
			maxNameLen = len(sourceFile.Name)
		}
		if len(sourceFile.Path) > maxPathLen {
			maxPathLen = len(sourceFile.Path)
		}
	}
	return maxNameLen, maxPathLen
}

func (l *List) renderSourceFiles(redraw bool, promptHeadroom int) {

	boldUnderlineAccent := style.CreateStyle().AddStyles(colors.TEXT_UNDERLINE, colors.TEXT_RED, colors.TEXT_BOLD)
	boldAccent := style.CreateStyle().AddStyles(colors.TEXT_RED)

	maxNameLen, maxPathLen := l.calculateMaxNameAndPathLen()

	if redraw {
		linesUp := len(l.SourceFiles) + promptHeadroom

		if linesUp > 0 {
			fmt.Printf("\033[%dA", linesUp+1)
		}

		for i := 0; i < linesUp; i++ {
			fmt.Printf("\033[2K\r\033[1B") // Clear line and move down
		}

		if linesUp > 0 {
			fmt.Printf("\033[%dA", linesUp+2)
		}
	}
	fmt.Printf("%s", l.Prompt)
	fmt.Print("\n")
	fmt.Printf("%s", "   ")
	fmt.Printf("%-15s", "Active")
	fmt.Printf("%-*s", maxNameLen+20, "File name")
	fmt.Printf("%-*s", maxPathLen+20, "File path")
	fmt.Printf("%-20s", "Priority")
	fmt.Print("\n")

	for i, SourceFile := range l.SourceFiles {
		cursor := "   "

		selection := "[ ]"

		if SourceFile.Selected {

			selection = "[" + boldAccent.Use("x") + "]"
		}

		if i == l.CursorPos {
			cursor = ">  "
		}

		preActiveSpacer := 15

		if SourceFile.Selected {
			preActiveSpacer += 9
		}

		prioDisplay := ""
		if SourceFile.Priority == "low" {
			prioDisplay = fmt.Sprintf("%s%s%s", "[", boldUnderlineAccent.Use("Low"), " / High]")
		} else if SourceFile.Priority == "high" {
			prioDisplay = fmt.Sprintf("%s%s%s", "[Low / ", boldUnderlineAccent.Use("High"), "]")
		}
		fmt.Printf("%s", cursor)
		fmt.Printf("%-*s", preActiveSpacer, selection)
		fmt.Printf("%-*s", maxNameLen+20, SourceFile.Name)
		fmt.Printf("%-*s", maxPathLen+20, SourceFile.Path)
		fmt.Printf("%-20s", prioDisplay)
		fmt.Print("\n")

	}
}

func (l *List) Display() (SourceFilePointerArray, bool) {
	defer func() {
		fmt.Printf("")
	}()

	l.renderSourceFiles(false, 0)

	for {

		currentFile := l.SourceFiles[l.CursorPos]

		keyCode := cli.CheckInput()

		// EXIT
		if keyCode == constants.Escape {
			return nil, false

			// ACCEPT
		} else if keyCode == constants.Enter {
			return l.SourceFiles, true

			// SELECT
		} else if keyCode == constants.Space {

			currentFile.Selected = !currentFile.Selected
			l.renderSourceFiles(true, 0)

			// ORDERING
		} else if keyCode == constants.MoveDown || keyCode == constants.MoveUp {

			if keyCode == constants.MoveDown {
				if l.CursorPos == 0 {
					l.SourceFiles = append(l.SourceFiles[1:], l.SourceFiles[0])
					l.CursorPos = len(l.SourceFiles) - 1
				} else {
					// Swap as before
					l.SourceFiles[l.CursorPos], l.SourceFiles[l.CursorPos-1] = l.SourceFiles[l.CursorPos-1], l.SourceFiles[l.CursorPos]
					l.CursorPos--
				}
			} else if keyCode == constants.MoveUp {
				if l.CursorPos == len(l.SourceFiles)-1 {
					l.SourceFiles = append([]*SourceFile{l.SourceFiles[len(l.SourceFiles)-1]}, l.SourceFiles[:len(l.SourceFiles)-1]...)
					l.CursorPos = 0
				} else {
					// Swap as before
					l.SourceFiles[l.CursorPos], l.SourceFiles[l.CursorPos+1] = l.SourceFiles[l.CursorPos+1], l.SourceFiles[l.CursorPos]
					l.CursorPos++
				}
			}

			l.renderSourceFiles(true, 0)

			// MOVEMENTS
		} else if keyCode == constants.Up || keyCode == constants.Down || keyCode == constants.Right || keyCode == constants.Left {

			if keyCode == constants.Up {
				l.CursorPos = (l.CursorPos + len(l.SourceFiles) - 1) % len(l.SourceFiles)
			}
			if keyCode == constants.Down {
				l.CursorPos = (l.CursorPos + 1) % len(l.SourceFiles)
			}
			if keyCode == constants.Left || keyCode == constants.Right {
				if strings.ToLower(currentFile.Priority) == "high" {
					currentFile.Priority = "low"
				} else {
					currentFile.Priority = "high"
				}
			}

			l.renderSourceFiles(true, 0)

			// RENAMING CURRENT SOURCE FILE
		} else if keyCode == constants.Rename {
			name := helper.GetUserString("Enter the name for the source file:")
			l.SourceFiles[l.CursorPos].Name = name
			l.renderSourceFiles(true, 3)

		} else if keyCode == constants.Create {
			path := helper.GetUserString("Enter the path for the source file:")
			name := helper.GetUserString("Enter the name for the source file:")

			fmt.Print("\033[H\033[2J")

			if !findIfCreatedP(l.SourceFiles, path) {
				l.AddSourceToList(SourceFile{
					Name:     name,
					Path:     path,
					Priority: "low",
					Selected: false,
				})
			} else {
				fmt.Print("File already exists.")
			}

			l.renderSourceFiles(true, 1)
		} else if keyCode == constants.Delete {
			answ := strings.ToLower(helper.GetUserString("Are you sure? [y/n]"))
			ok := answ == "yes" || answ == "y"

			if ok {
				l.SourceFiles = append(l.SourceFiles[:l.CursorPos], l.SourceFiles[l.CursorPos+1:]...)
				l.CursorPos = 0
			}
			fmt.Print("\033[H\033[2J")
			l.renderSourceFiles(true, 4)
		}
	}
}

func findIfCreatedP(jsonSources []*SourceFile, source string) bool {
	for _, jsonSource := range jsonSources {
		if strings.TrimSpace(source) == jsonSource.Path {
			return true
		}
	}
	return false
}
