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
	NameToggled bool
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

func (l *List) renderSourceFiles() {

	boldUnderlineAccent := style.CreateStyle().AddStyles(colors.TEXT_UNDERLINE, colors.TEXT_RED, colors.TEXT_BOLD)
	boldAccent := style.CreateStyle().AddStyles(colors.TEXT_RED)
	fmt.Print(colors.CLEAR_SCREEN)

	fmt.Printf("%s", l.Prompt)
	fmt.Print("\n")
	fmt.Printf("%s", "   ")
	fmt.Printf("%-15s", "Active")
	if !l.NameToggled {
		fmt.Printf("%-*s", 40, "File name")
	} else {
		fmt.Printf("%-*s", 40, "File path")
	}
	fmt.Printf("%-5s", "Priority")
	fmt.Print("\n\n")

	if len(l.SourceFiles) > 0 {

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

			filename := ""
			filepath := ""

			for i, char := range SourceFile.Name {
				if i < 35 {
					filename += string(char)
				} else {
					filename += "..."
					break
				}
			}

			for i, char := range SourceFile.Path {
				if i < 35 {
					filepath += string(char)
				} else {
					filepath += "..."
					break
				}
			}

			if !l.NameToggled {
				fmt.Printf("%-*s", 40, filename)
			} else {
				fmt.Printf("%-*s", 40, filepath)
			}

			fmt.Printf("%-5s", prioDisplay)
			fmt.Print("\n")

		}

	} else {
		fmt.Printf("You don't have any source files, to create one, press [a]\n")
	}

}

func (l *List) Display() (SourceFilePointerArray, bool) {
	defer func() {
		fmt.Printf("")
	}()

	boldUnderlineAccent := style.CreateStyle().AddStyles(colors.TEXT_UNDERLINE, colors.TEXT_RED, colors.TEXT_BOLD)
	redNameString := fmt.Sprintf(boldUnderlineAccent.Use("name"))
	redPathString := fmt.Sprintf(boldUnderlineAccent.Use("path"))

	l.renderSourceFiles()

	for {

		keyCode := cli.CheckInput()

		if len(l.SourceFiles) < 1 {

			if keyCode == constants.Escape {
				fmt.Print(colors.CLEAR_SCREEN)
				return nil, false

				// ACCEPT
			} else if keyCode == constants.Enter {
				fmt.Print(colors.CLEAR_SCREEN)
				return l.SourceFiles, true

				// SELECT
			} else if keyCode == constants.Create {

				var name string
				var path string
				ok := true

				path = helper.GetUserString("Enter the " + redPathString + " for the source file (empty to cancel): ")
				if path == "" {
					ok = false
				}
				if ok {
					name = helper.GetUserString("Enter the " + redNameString + " for the source file (empty to cancel): ")
					if name == "" {
						ok = false
					}
				}
				if ok {
					l.AddSourceToList(SourceFile{
						Name:     name,
						Path:     path,
						Priority: "low",
						Selected: false,
					})
				}

				l.renderSourceFiles()
			}

		} else {

			currentFile := l.SourceFiles[l.CursorPos]

			// EXIT
			if keyCode == constants.Escape {
				fmt.Print(colors.CLEAR_SCREEN)
				return nil, false

				// ACCEPT
			} else if keyCode == constants.Enter {
				fmt.Print(colors.CLEAR_SCREEN)
				return l.SourceFiles, true

				// SELECT
			} else if keyCode == constants.Space {

				currentFile.Selected = !currentFile.Selected
				l.renderSourceFiles()

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

				l.renderSourceFiles()

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

				l.renderSourceFiles()

				// RENAMING CURRENT SOURCE FILE
			} else if keyCode == constants.Rename {
				name := helper.GetUserString("Enter the name for the source file (empty to cancel):")
				if name != "" {
					l.SourceFiles[l.CursorPos].Name = name
				}
				l.renderSourceFiles()

				// CREATING A SOURCE FILE
			} else if keyCode == constants.Create {

				var name string
				var path string
				ok := true

				path = helper.GetUserString("Enter the " + redPathString + " for the source file (empty to cancel): ")
				if path == "" {
					ok = false
				}
				if ok {
					name = helper.GetUserString("Enter the " + redNameString + " for the source file (empty to cancel): ")
					if name == "" {
						ok = false
					}
				}
				if ok && !findIfCreatedP(l.SourceFiles, path) {
					l.AddSourceToList(SourceFile{
						Name:     name,
						Path:     path,
						Priority: "low",
						Selected: false,
					})
				} else {
					fmt.Print("File already exists.")
				}

				l.renderSourceFiles()

			} else if keyCode == constants.Delete {
				answ := strings.ToLower(helper.GetUserString("Are you sure? [y/n] "))
				ok := answ == "yes" || answ == "y"

				if ok {
					l.SourceFiles = append(l.SourceFiles[:l.CursorPos], l.SourceFiles[l.CursorPos+1:]...)
					l.CursorPos = 0
				}
				l.renderSourceFiles()

			} else if keyCode == constants.ToggleName {
				l.NameToggled = !l.NameToggled
				l.renderSourceFiles()

			} else if keyCode == constants.Help {
				fmt.Print(colors.CLEAR_SCREEN)
				fmt.Print("Hyprsources is a tool to quickly manage your Hyprland source files.\n")
				fmt.Print("If you experience problems, don't hesitate to create an issue or a PR in Github:\nhttps://github.com/andrwui/hyprsources\n")
				fmt.Println()
				fmt.Println("=== CONTROLS ===")
				fmt.Println("[return] - Exit and save changes")
				fmt.Println("[escape] - Exit without saving changes")
				fmt.Println("[space] - Select source as active")
				fmt.Println("[up/down] - Move the cursor")
				fmt.Println("[right/left] - Change priority")
				fmt.Println("[j/k] - Change source load order")
				fmt.Println("[a] - Create a source")
				fmt.Println("[r] - Rename a source")
				fmt.Println("[d] - Delete a source")
				fmt.Println("[r] - Rename a source")
				fmt.Println("[l] - Toggle between file path and source name")
				fmt.Println("Press any key to go back")
				cli.CheckInput()
				l.renderSourceFiles()
			}
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
