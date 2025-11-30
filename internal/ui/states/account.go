// account settings view for the ui to login / register and manage account preferences
package states

import (
	"strings"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"cliscraper/internal/ui/model"
	"cliscraper/internal/ui/components"
	"cliscraper/internal/ui/messages"
)

type AuthDoneMsg = messages.AuthDoneMsg

func UpdateAccount(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			// switch between login and register modes
			if m.AccountMode == "login" {
				m.AccountMode = "register"
				m.AccountField = "username"
				m.Username = ""
				m.Password = ""
				m.Email = ""
				m.Err = ""
			} else {
				m.AccountMode = "login"
				m.AccountField = "username"
				m.Username = ""
				m.Password = ""
				m.Email = ""
				m.Err = ""
			}
			return m, nil
		case tea.KeyEnter:
			// submission
			if m.AccountMode == "login" {
				if m.Username == "" || m.Password == "" {
					m.Err = "Username and password are required"
					return m, nil
				}
				// start login async operation
				m.Spinner = components.InitialSpinner()
				m.Err = ""
				return m, tea.Batch(
					m.Spinner.Init(),
					LoginCmd(m, m.Username, m.Password),
				)
			} else { // register
				if m.Username == "" || m.Email == "" || m.Password == "" {
					m.Err = "Username, email, and password are required"
					return m, nil
				}
				if len(m.Password) < 6 {
					m.Err = "Password must be at least 6 characters"
					return m, nil
				}
				// start register async operation
				m.Spinner = components.InitialSpinner()
				m.Err = ""
				return m, tea.Batch(
					m.Spinner.Init(),
					RegisterCmd(m, m.Username, m.Email, m.Password),
				)
			}
		case tea.KeyBackspace, tea.KeyDelete:
			// handle backspace based on current field
			switch m.AccountField {
			case "username":
				if len(m.Username) > 0 {
					m.Username = m.Username[:len(m.Username)-1]
				}
			case "email":
				if len(m.Email) > 0 {
					m.Email = m.Email[:len(m.Email)-1]
				}
			case "password":
				if len(m.Password) > 0 {
					m.Password = m.Password[:len(m.Password)-1]
				}
			}
			return m, nil
		case tea.KeyDown, tea.KeyUp:
			// move to next/previous field
			if m.AccountMode == "login" {
				if m.AccountField == "username" {
					m.AccountField = "password"
				} else {
					m.AccountField = "username"
				}
			} else { // register
				if msg.Type == tea.KeyDown {
					if m.AccountField == "username" {
						m.AccountField = "email"
					} else if m.AccountField == "email" {
						m.AccountField = "password"
					} else {
						m.AccountField = "username"
					}
				} else { // key up
					if m.AccountField == "username" {
						m.AccountField = "password"
					} else if m.AccountField == "password" {
						m.AccountField = "email"
					} else {
						m.AccountField = "username"
					}
				}
			}
			return m, nil
		default:
			// handle text input
			if msg.Type == tea.KeyRunes {
				char := msg.String()
				switch m.AccountField {
				case "username":
					m.Username += char
				case "email":
					m.Email += char
				case "password":
					m.Password += char
				}
			}
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	case messages.AuthDoneMsg:
		if msg.Err != nil {
			m.Err = msg.Err.Error()
		} else {
			m.CurrentUser = msg.User
			m.Err = ""
			// clear
			m.Username = ""
			m.Email = ""
			m.Password = ""
			m.AccountField = "username"
			// return to home after successful auth
			m.CurrentState = model.StateHome
		}
		return m, nil
	}
	return m, nil
}

func ViewAccountSettings(m model.Model) string {
	var b strings.Builder

	// title
	modeTitle := "Login"
	if m.AccountMode == "register" {
		modeTitle = "Register"
	}
	b.WriteString(components.TitleStyle.Render(fmt.Sprintf("Account Settings - %s", modeTitle)) + "\n\n")

	// mode switcher hint
	modeHint := "Press TAB to switch to Register"
	if m.AccountMode == "register" {
		modeHint = "Press TAB to switch to Login"
	}
	b.WriteString(components.StatusStyle.Render(modeHint) + "\n\n")

	// show spinner if authenticating
	spinnerView := m.Spinner.View()
	if spinnerView != "" {
		b.WriteString(components.StatusStyle.Render(fmt.Sprintf("%s Authenticating...", spinnerView)) + "\n\n")
	}

	// form fields
	if m.AccountMode == "login" {
		// username field
		usernameLabel := "Username: "
		if m.AccountField == "username" {
			usernameLabel = "> Username: "
		}
		b.WriteString(components.LabelStyle.Render(usernameLabel))
		b.WriteString(components.InputStyle.Render(m.Username) + "\n")

		// password field
		passwordLabel := "Password: "
		if m.AccountField == "password" {
			passwordLabel = "> Password: "
		}
		b.WriteString(components.LabelStyle.Render(passwordLabel))
		// mask password, for screen sharing safety etc
		maskedPassword := strings.Repeat("*", len(m.Password))
		b.WriteString(components.InputStyle.Render(maskedPassword) + "\n")
	} else { // register
		// username field
		usernameLabel := "Username: "
		if m.AccountField == "username" {
			usernameLabel = "> Username: "
		}
		b.WriteString(components.LabelStyle.Render(usernameLabel))
		b.WriteString(components.InputStyle.Render(m.Username) + "\n")

		// email field
		emailLabel := "Email: "
		if m.AccountField == "email" {
			emailLabel = "> Email: "
		}
		b.WriteString(components.LabelStyle.Render(emailLabel))
		b.WriteString(components.InputStyle.Render(m.Email) + "\n")

		// password field
		passwordLabel := "Password: "
		if m.AccountField == "password" {
			passwordLabel = "> Password: "
		}
		b.WriteString(components.LabelStyle.Render(passwordLabel))
		// mask password again
		maskedPassword := strings.Repeat("*", len(m.Password))
		b.WriteString(components.InputStyle.Render(maskedPassword) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(components.StatusStyle.Render("Press ENTER to submit, TAB to switch mode, ↑/↓ to switch fields") + "\n")

	// show current user on ui if logged in
	if m.CurrentUser != nil {
		b.WriteString("\n")
		b.WriteString(components.StatusStyle.Render(fmt.Sprintf("Logged in as: %s (%s)", m.CurrentUser.Username, m.CurrentUser.Email)) + "\n")
	}

	return b.String()
}

// perform async login
func LoginCmd(m model.Model, username, password string) tea.Cmd {
	return func() tea.Msg {
		userInfo, err := m.Service().Login(username, password)
		if err != nil {
			return AuthDoneMsg{Err: err}
		}
		return AuthDoneMsg{User: userInfo}
	}
}

// perform async registration
func RegisterCmd(m model.Model, username, email, password string) tea.Cmd {
	return func() tea.Msg {
		userInfo, err := m.Service().Register(username, email, password)
		if err != nil {
			return AuthDoneMsg{Err: err}
		}
		return AuthDoneMsg{User: userInfo}
	}
}
