package simpleform

import (
	"encoding/json"
	"fmt"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
)

// Struct ==================================================================
type SimpleForm struct {
	title   string
	desc    string
	btns    []buttonData
	onClose func(p *player.Player)
}

type buttonData struct {
	btn     form.Button
	onClick func(p *player.Player)
}

type menuHandler struct{ *SimpleForm }

type iconObject struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type button struct {
	Text string      `json:"text"`
	Icon *iconObject `json:"image,omitempty"`
}

// Struct ==================================================================
// Functions ==================================================================
func (m *menuHandler) Title() string { return m.title }
func (m *menuHandler) Body() string  { return m.desc }
func (m *menuHandler) Buttons() []form.Button {
	btns := make([]form.Button, len(m.btns))
	for i, b := range m.btns {
		btns[i] = b.btn
	}
	return btns
}

func (m *menuHandler) Submit(p *player.Player, idx int) {
	if idx >= 0 && idx < len(m.btns) {
		if m.btns[idx].onClick != nil {
			m.btns[idx].onClick(p)
		}
	}
}

func (m *menuHandler) Close(p *player.Player) {
	if m.onClose != nil {
		m.onClose(p)
	}
}

func (m *menuHandler) MarshalJSON() ([]byte, error) {
	btns := make([]button, len(m.btns))
	for i, b := range m.btns {
		var icon *iconObject
		if b.btn.Image != "" {
			icon = &iconObject{
				Type: "path",
				Data: b.btn.Image,
			}
		}
		btns[i] = button{Text: b.btn.Text, Icon: icon}
	}
	data := map[string]any{
		"type":    "form",
		"title":   m.title,
		"content": m.desc,
		"buttons": btns,
	}
	return json.Marshal(data)
}

func (m *menuHandler) SubmitJSON(b []byte, submitter form.Submitter, tx *world.Tx) error {
	if len(b) == 0 || string(b) == "null" {
		if m.onClose != nil {
			if p, ok := submitter.(*player.Player); ok {
				m.onClose(p)
			}
		}
		return nil
	}
	var idx int
	if err := json.Unmarshal(b, &idx); err != nil {
		return fmt.Errorf("gagal parse submit index: %w", err)
	}
	if p, ok := submitter.(*player.Player); ok {
		m.Submit(p, idx)
	}
	return nil
}

// Functions ==================================================================
// Builder ==================================================================
// Title & Description
func New(title, desc string) *SimpleForm { return &SimpleForm{title: title, desc: desc} }

// B: Add buttons
func (f *SimpleForm) B(text, path string, cb func(p *player.Player)) *SimpleForm {
	f.btns = append(f.btns, buttonData{
		btn:     form.NewButton(text, path),
		onClick: cb,
	})
	return f
}

// Close callback
func (f *SimpleForm) Close(cb func(p *player.Player)) *SimpleForm { f.onClose = cb; return f }

// S: Send form to player
func (f *SimpleForm) S(p *player.Player) { p.SendForm(&menuHandler{f}) }

// Builder ==================================================================
