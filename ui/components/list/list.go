package list

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Item represents a standard list item with title and description.
// It implements list.DefaultItem and list.Item.
type Item struct {
	title       string
	description string
}

// NewItem creates a new Item with a title and description.
func NewItem(title, description string) Item {
	return Item{
		title:       title,
		description: description,
	}
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.description }
func (i Item) FilterValue() string { return i.title }

type ListModel struct {
	List         list.Model
	SelectedItem list.Item
	Done         bool
	Quitting     bool
	Err          error
}

// CreateListModel initializes a ListModel with default dimensions and a default delegate.
func CreateListModel(items []list.Item, title string) ListModel {
	return CreateListModelWithDimensions(items, title, 40, 14)
}

// CreateListModelWithDimensions initializes a ListModel with specified dimensions.
func CreateListModelWithDimensions(items []list.Item, title string, width, height int) ListModel {
	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, width, height)
	l.Title = title
	l.SetShowStatusBar(false)

	return ListModel{
		List: l,
	}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		if m.List.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit

		case "enter":
			if selected := m.List.SelectedItem(); selected != nil {
				m.SelectedItem = selected
				m.Done = true
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	if m.Quitting {
		return ""
	}
	return m.List.View()
}

// Value returns the title of the selected item, or empty string if nothing is selected.
func (m ListModel) Value() string {
	if m.SelectedItem == nil {
		return ""
	}
	if defaultItem, ok := m.SelectedItem.(list.DefaultItem); ok {
		return defaultItem.Title()
	}
	return m.SelectedItem.FilterValue()
}

// Index returns the index of the currently highlighted item in the list.
func (m ListModel) Index() int {
	return m.List.Index()
}

