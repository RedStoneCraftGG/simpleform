package simpleform

import (
	"encoding/json"
	"fmt"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
)

// ===================== Structs =====================

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

type SubmitForm struct {
	title    string
	desc     string
	elements []submitElement
	onSubmit func(p *player.Player, r SubmitFormResponse)
	onClose  func(p *player.Player)
}

type submitElement struct {
	typ            string
	label          string
	options        []string
	defaultInt     int
	defaultBool    bool
	defaultFloat   int
	placeholder    string
	defaultText    string
	min, max, step int
}

type SubmitFormResponse struct {
	values   []any
	elements []submitElement
}

type submitFormHandler struct{ *SubmitForm }

type ModalForm struct {
	title    string
	content  string
	btn1Text string
	btn2Text string
	btn1Func func(p *player.Player)
	btn2Func func(p *player.Player)
	onClose  func(p *player.Player)
}

type modalFormHandler struct{ *ModalForm }

// ===================== SimpleForm Methods =====================

func (m *menuHandler) Title() string { return m.title }

func (m *menuHandler) Body() string { return m.desc }

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
		return fmt.Errorf("failed to parse submit index: %w", err)
	}
	if p, ok := submitter.(*player.Player); ok {
		m.Submit(p, idx)
	}
	return nil
}

// ===================== SimpleForm Builder =====================

func New(title, desc string) *SimpleForm { return &SimpleForm{title: title, desc: desc} }

func (f *SimpleForm) B(text, path string, cb func(p *player.Player)) *SimpleForm {
	f.btns = append(f.btns, buttonData{
		btn:     form.NewButton(text, path),
		onClick: cb,
	})
	return f
}

func (f *SimpleForm) Close(cb func(p *player.Player)) *SimpleForm {
	f.onClose = cb
	return f
}

func (f *SimpleForm) S(p *player.Player) { p.SendForm(&menuHandler{f}) }

// ===================== SubmitForm Builder =====================

func NewSubmitForm(title, desc string) *SubmitForm { return &SubmitForm{title: title, desc: desc} }

func (f *SubmitForm) Dropdown(label string, options []string, defaultIndex ...int) *SubmitForm {
	idx := 0
	if len(defaultIndex) > 0 {
		idx = defaultIndex[0]
	}
	f.elements = append(f.elements, submitElement{
		typ:        "dropdown",
		label:      label,
		options:    options,
		defaultInt: idx,
	})
	return f
}

func (f *SubmitForm) Toggle(label string, defaultValue ...bool) *SubmitForm {
	val := false
	if len(defaultValue) > 0 {
		val = defaultValue[0]
	}
	f.elements = append(f.elements, submitElement{
		typ:         "toggle",
		label:       label,
		defaultBool: val,
	})
	return f
}

func (f *SubmitForm) Slider(label string, min, max, step int, defaultValue ...int) *SubmitForm {
	val := min
	if len(defaultValue) > 0 {
		val = defaultValue[0]
	}
	f.elements = append(f.elements, submitElement{
		typ:          "slider",
		label:        label,
		min:          min,
		max:          max,
		step:         step,
		defaultFloat: val,
	})
	return f
}

func (f *SubmitForm) Input(label, placeholder string, defaultValue ...string) *SubmitForm {
	val := ""
	if len(defaultValue) > 0 {
		val = defaultValue[0]
	}
	f.elements = append(f.elements, submitElement{
		typ:         "input",
		label:       label,
		placeholder: placeholder,
		defaultText: val,
	})
	return f
}

func (f *SubmitForm) OnSubmit(cb func(p *player.Player, r SubmitFormResponse)) *SubmitForm {
	f.onSubmit = cb
	return f
}

func (f *SubmitForm) OnClose(cb func(p *player.Player)) *SubmitForm {
	f.onClose = cb
	return f
}

func (f *SubmitForm) S(p *player.Player) { p.SendForm(&submitFormHandler{f}) }

// ===================== SubmitFormResponse Methods =====================

func (r SubmitFormResponse) Dropdown(idx int) string {
	if idx < len(r.values) && idx < len(r.elements) {
		opt := toInt(r.values[idx])
		elem := r.elements[idx]
		if elem.typ == "dropdown" && opt >= 0 && opt < len(elem.options) {
			return elem.options[opt]
		}
	}
	return ""
}

func (r SubmitFormResponse) DropdownIndex(idx int) int {
	if idx < len(r.values) {
		return toInt(r.values[idx])
	}
	return 0
}

func (r SubmitFormResponse) Toggle(idx int) bool {
	if idx < len(r.values) {
		return toBool(r.values[idx])
	}
	return false
}

func (r SubmitFormResponse) Slider(idx int) int {
	if idx < len(r.values) {
		return toInt(r.values[idx])
	}
	return 0
}

func (r SubmitFormResponse) Input(idx int) string {
	if idx < len(r.values) {
		return toString(r.values[idx])
	}
	return ""
}

func (r SubmitFormResponse) valuesDropdownOptions(idx int) []string {
	if idx < len(r.elements) && r.elements[idx].typ == "dropdown" {
		return r.elements[idx].options
	}
	return nil
}

// ===================== Helpers =====================

func toInt(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case int32:
		return int(val)
	case int64:
		return int(val)
	case float64:
		return int(val)
	case float32:
		return int(val)
	default:
		return 0
	}
}

func toBool(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	default:
		return false
	}
}

func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		return ""
	}
}

// ===================== submitFormHandler Methods =====================

func (h *submitFormHandler) Title() string { return h.title }

func (h *submitFormHandler) Body() string { return h.desc }

func (h *submitFormHandler) Buttons() []form.Button { return nil }

func (h *submitFormHandler) MarshalJSON() ([]byte, error) {
	var content []any
	for _, e := range h.elements {
		switch e.typ {
		case "dropdown":
			m := map[string]any{
				"type":    "dropdown",
				"text":    e.label,
				"options": e.options,
			}
			if e.defaultInt != 0 {
				m["default"] = e.defaultInt
			}
			content = append(content, m)
		case "toggle":
			m := map[string]any{
				"type": "toggle",
				"text": e.label,
			}
			if e.defaultBool {
				m["default"] = e.defaultBool
			}
			content = append(content, m)
		case "slider":
			m := map[string]any{
				"type": "slider",
				"text": e.label,
				"min":  e.min,
				"max":  e.max,
				"step": e.step,
			}
			if e.defaultFloat != 0 {
				m["default"] = e.defaultFloat
			}
			content = append(content, m)
		case "input":
			m := map[string]any{
				"type":        "input",
				"text":        e.label,
				"placeholder": e.placeholder,
			}
			if e.defaultText != "" {
				m["default"] = e.defaultText
			}
			content = append(content, m)
		}
	}
	data := map[string]any{
		"type":    "custom_form",
		"title":   h.title,
		"content": content,
	}
	return json.Marshal(data)
}

func (h *submitFormHandler) SubmitJSON(b []byte, submitter form.Submitter, tx *world.Tx) error {
	if len(b) == 0 || string(b) == "null" {
		if h.onClose != nil {
			if p, ok := submitter.(*player.Player); ok {
				h.onClose(p)
			}
		}
		return nil
	}

	var arr []any
	if err := json.Unmarshal(b, &arr); err != nil {
		return fmt.Errorf("failed to parse submit: %w", err)
	}

	resp := SubmitFormResponse{values: arr, elements: h.elements}
	if h.onSubmit != nil {
		if p, ok := submitter.(*player.Player); ok {
			h.onSubmit(p, resp)
		}
	}

	return nil
}

func (h *submitFormHandler) Submit(p *player.Player, idx int) {}

func (h *submitFormHandler) Close(p *player.Player) {
	if h.onClose != nil {
		h.onClose(p)
	}
}

// ===================== ModalForm Structs & Methods =====================

func NewModalForm(title, content string) *ModalForm {
	return &ModalForm{
		title:   title,
		content: content,
	}
}

func (m *ModalForm) B1(text string, f func(p *player.Player)) *ModalForm {
	m.btn1Text = text
	m.btn1Func = f
	return m
}

func (m *ModalForm) B2(text string, f func(p *player.Player)) *ModalForm {
	m.btn2Text = text
	m.btn2Func = f
	return m
}

func (m *ModalForm) Close(f func(p *player.Player)) *ModalForm {
	m.onClose = f
	return m
}

func (m *ModalForm) S(p *player.Player) { p.SendForm(&modalFormHandler{m}) }

// ===================== modalFormHandler Implementation =====================

func (h *modalFormHandler) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"type":    "modal",
		"title":   h.title,
		"content": h.content,
		"button1": h.btn1Text,
		"button2": h.btn2Text,
	}
	return json.Marshal(data)
}

func (h *modalFormHandler) SubmitJSON(b []byte, submitter form.Submitter, tx *world.Tx) error {
	if len(b) == 0 || string(b) == "null" {
		if h.onClose != nil {
			if p, ok := submitter.(*player.Player); ok {
				h.onClose(p)
			}
		}
		return nil
	}

	var btnIdx bool
	if err := json.Unmarshal(b, &btnIdx); err != nil {
		return fmt.Errorf("failed to parse modal submit: %w", err)
	}

	if p, ok := submitter.(*player.Player); ok {
		if btnIdx {
			if h.btn1Func != nil {
				h.btn1Func(p)
			}
		} else {
			if h.btn2Func != nil {
				h.btn2Func(p)
			}
		}
	}
	return nil
}

func (h *modalFormHandler) Submit(p *player.Player, idx int) {}

func (h *modalFormHandler) Close(p *player.Player) {
	if h.onClose != nil {
		h.onClose(p)
	}
}
