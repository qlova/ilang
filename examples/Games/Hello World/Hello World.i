import grate

type Graphics {}

method draw(Graphics) {
	set.color(255, 255, 255)
	draw(Text{"Hello World"})
	set.color(0, 0, 0)
}
