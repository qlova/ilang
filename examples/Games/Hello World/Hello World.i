import grate

type Game {}

method draw(Game) {
	set(Color{255, 255, 0})
	draw(Text{"Hello World"})
	set(Color{0, 0, 0})
}
