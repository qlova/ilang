import grate

const gravity = 10

type Cube {x, y, jumptimer, jumping}

method update(Cube) {
	if key.a()
		x -= 10px
	end
	if key.d()
		x += 10px
	end
	
	x += gamepad(0).axis(0)
	
	if (key.space() + gamepad(0).down(2)) * (jumptimer > 0)
		jumptimer--
		y -= gravity px 
		y -= (jumptimer) px
	end
	
	
	y += gravity px
	if y >= (height() - 50px)
		jumptimer = 10
		y = height() - 50px
	end
}

method draw(Cube) {
	square(x, y, 100px)
}

type Graphics {
	{Cube}Cube
}

method new(Graphics) {
	Cube.x = 200px
	Cube.y = 200px
}

method update(Graphics) {
	update(Cube)
}

method draw(Graphics) {
	draw(Cube)
}
