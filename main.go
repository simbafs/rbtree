package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	tree *RBTree
	cmd  textinput.Model
	msg  string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.cmd.SetValue("")
		case "enter":
			if m.cmd.Value() == "" {
				break
			}
			// c -> check
			// o - print in order
			// r 10  -> right rotate node 10
			// l 50 -> left rotate node 50
			// p 30  -> print node 30
			// r 30.left.right.parent -> print 30's left's right's parent
			args := strings.Split(m.cmd.Value(), " ")
			switch args[0] {
			case "c":
				// m.msg = "not implemented yet"
				m.check()
			case "o":
				m.msg = m.tree.InOrder()
			case "r":
				if len(args) < 2 {
					m.msg = "Please provide a value to rotate"
					break
				}

				node, err := m.tree.Query(args[1])
				if err != nil {
					m.msg = err.Error()
					break
				}

				m.tree.RRotate(node)

			case "l":
				if len(args) < 2 {
					m.msg = "Please provide a value to rotate"
					break
				}

				node, err := m.tree.Query(args[1])
				if err != nil {
					m.msg = err.Error()
					break
				}

				m.tree.LRotate(node)
			case "p":
				if len(args) < 2 {
					m.print("")
				} else {
					m.print(args[1])
				}
			case "i":
				if len(args) < 2 {
					m.msg = "Please provide a value to insert"
					break
				}
				val, err := strconv.Atoi(args[1])
				if err != nil {
					m.msg = "Invalid value " + args[1]
					break
				}

				m.tree.Insert(val)
			case "d":
				if len(args) < 2 {
					m.msg = "Please provide a value to delete"
					break
				}
				val, err := strconv.Atoi(args[1])
				if err != nil {
					m.msg = "Invalid value " + args[1]
					break
				}

				m.tree.Delete(m.tree.Find(val))
			case "q":
				return m, tea.Quit
			default:
				m.msg = "Unknown command " + args[0]
			}

			m.cmd.SetValue("")
		}
	}

	m.cmd, cmd = m.cmd.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		"Trees:",
		m.tree.View(),
		m.cmd.View(),
		lipgloss.JoinHorizontal(lipgloss.Left,
			"< ",
			m.msg,
		),
	)
}

func (m *Model) check() {
	queue := []*RBNode{m.tree.Root}
	m.msg = ""
	if m.tree.Root.Parent != Nil {
		m.msg += "Root's parent is not nil\n"
	}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == nil || curr == Nil {
			continue
		}

		m.msg += fmt.Sprintf("Node %d: \n", curr.Value)

		if curr.Left != Nil && curr.Left.Parent != curr {
			m.msg += fmt.Sprintf("%d's parent is not %d but %d\n", curr.Left.Value, curr.Value, curr.Left.Parent.Value)
		}

		if curr.Right != Nil && curr.Right.Parent != curr {
			m.msg += fmt.Sprintf("%d's parent is not %d but %d\n", curr.Right.Value, curr.Value, curr.Right.Parent.Value)
		}

		queue = append(queue, curr.Left, curr.Right)
	}
}

func (m *Model) print(query string) {
	m.msg = ""
	if query == "" {
		// print all
		queue := []*RBNode{m.tree.Root}
		m.msg = ""
		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]

			if curr == nil {
				continue
			}

			m.msg += curr.String() + "\n"

			queue = append(queue, curr.Left, curr.Right)
		}

		return
	}

	paths := strings.Split(query, ".")
	curr := m.tree.Root

	for index, path := range paths {
		if curr == nil {
			m.msg += "Node is nil at path: " + strings.Join(paths[:index], ".") + "\n"
			continue
		}
		switch path {
		case "left":
			curr = curr.Left
		case "right":
			curr = curr.Right
		case "parent":
			curr = curr.Parent
		case "pre":
			curr = curr.Predecessor()
		case "suc":
			curr = curr.Successor()
		case "root", "":
			curr = m.tree.Root
		default:
			if val, err := strconv.Atoi(path); err == nil {
				curr = m.tree.Find(val)
			} else {
				m.msg = "Unknown path " + query
			}

		}
	}

	m.msg += curr.String()
}

func initModel() Model {
	cmd := textinput.New()
	// cmd.Placeholder = "> "
	cmd.Focus()

	return Model{
		tree: NewRBTree(),
		cmd:  cmd,
	}
}

func main() {
	tea.LogToFile("log.txt", "log")
	model := initModel()

	for _, v := range []int{50, 10, 80, 65, 90, 60, 70}{
		model.tree.Insert(v)
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
