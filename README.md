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

I couldn't figure out how to add a header, label, and divider. So, I'll leave some code unused inside. I hope it won't crash the code...

### Submit Form

```go
func TestSubmit(p *player.Player) {
	simpleform.NewSubmitForm("Settings", "Description (Not used, try using Label instead)").
		// Every element is counted as 1 index

		Header("This is a header"). // index 0 - Header Text
		Label("This is a label").   // index 1 - Label Text
		Divider().                  // index 2 - Divider

		// text: string, options: []string, default(optional): int
		Dropdown("Your favorite color:", []string{"Yellow", "Red", "Blue", "Green"}). // index 3 - Dropdown

		// text: string, default(optional): bool
		Toggle("Turn on night mode?"). // index 4 - Toggle

		// text: string, min: int, max: int, step: int, default(optional): int
		Slider("Volume", 0, 100, 5). // index 5 - Slider

		// text: string, placeholder: string, default(optional): string
		Input("Type your Nickname:", "Nickname"). // index 6 - Input

		OnSubmit(func(p *player.Player, r simpleform.SubmitFormResponse) {
			color := r.Dropdown(3)   // index 3
			nightMode := r.Toggle(4) // index 4
			volume := r.Slider(5)    // index 5
			nickname := r.Input(6)   // index 6
			p.Messagef("color: %s, Night Mode: %t, Volume: %v, Nick: %s", color, nightMode, volume, nickname)
		}).
		OnClose(func(p *player.Player) { /*optional*/ }).
		S(p)
}
```

<img width="768" height="768" alt="image" src="https://github.com/user-attachments/assets/077ee25e-cc7a-4f0a-ba98-87201b3d76e0" />


### Modal Form

```go
func TestModal(p *player.Player) {
	simpleform.NewModalForm("Conrim", "Do you want to have 271T added to your account?").
		B1("YESSSS PLEASE!!!", func(p *player.Player) { p.Message("271T was added to your account") }).
		B2("(Calmly Reject)", func(p *player.Player) { p.Message("You reject the offer") }).
		S(p)
}
```

<img width="768" height="768" alt="image" src="https://github.com/user-attachments/assets/ad4bf143-025b-447b-b105-3a8fda636f27" />
