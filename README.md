# simpleform
A lightweight and intuitive form builder for Dragonfly-MC servers, making it easy to create and manage button-based menus.

## Usage

### Normal Form

```go
func TestForm(p *player.Player) {
	sf.New("title", "desc").
		B("text1", "", func(p *player.Player) { p.Message("Clicked!") }).
		B("text2", "", func(p *player.Player) { p.Message("Clicked!") }).
		B("text3", "textures/red/carlotta", func(p *player.Player) { p.Message("Clicked!") }).
		Close(func(p *player.Player) { /*optional*/ }).
		S(p)
}
```

<img width="768" height="650" alt="image" src="https://github.com/user-attachments/assets/e0c29a0e-5478-4b08-96b0-6e60627200cd" />

## Still on update!

I'll add another form like SubmitForm and ModalForm later
