# simpleform
A lightweight and intuitive form builder for Dragonfly-MC servers, making it easy to create and manage button-based menus.

## Usage

### Normal Form

```go
func TestForm(p *player.Player) {
	simpleform.New("title", "desc").
		B("text1", "", func(p *player.Player) { p.Message("Clicked!") }).
		B("text2", "", func(p *player.Player) { p.Message("Clicked!") }).
		B("text3", "textures/red/carlotta", func(p *player.Player) { p.Message("Clicked!") }).
		Close(func(p *player.Player) { /*optional*/ }).
		S(p)
}
```

<img width="768" height="768" alt="image" src="https://github.com/user-attachments/assets/e0c29a0e-5478-4b08-96b0-6e60627200cd" />

### Submit Form

```go
func TestSubmit(p *player.Player) {
	simpleform.NewSubmitForm("Settings", "Description (Currently doesn't showed up)").

		// text: string, options: []string, default(optional): int
		Dropdown("Your favorite color:", []string{"Yellow", "Red", "Blue", "Green"}).

		// text: string, default(optional): bool
		Toggle("Turn on night mode?").

		// text: string, min: int, max: int, step: int, default(optional): int
		Slider("Volume:", 0, 100, 5).

		// text: string, placeholder: string, default(optional): string
		Input("Type your Nickname:", "Nickname").
		OnSubmit(func(p *player.Player, r simpleform.SubmitFormResponse) {
			color := r.Dropdown(0)   // index 0
			nightMode := r.Toggle(1) // index 1
			volume := r.Slider(2)    // index 2
			nickname := r.Input(3)   // index 3
			p.Messagef("color: %s, Night Mode: %t, Volume: %v, Nick: %s", color, nightMode, volume, nickname)
		}).
		OnClose(func(p *player.Player) { /*optional*/ }).
		S(p)
}
```

<img width="768" height="768" alt="image" src="https://github.com/user-attachments/assets/2ca6b238-c61e-485c-b486-20bffc13f2db" />


## Still on update!

I'll add another form like SubmitForm and ModalForm later
