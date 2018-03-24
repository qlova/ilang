import grate

function distance(x, y, a, b) {
	return (x-a)² + (y-b)²
}

type Graphics { angle, icon }

method new(Graphics) {
	icon = image.load("icon.png")
}

method update(Graphics) {
	angle += 1
	
	
}

method draw(Graphics) {
	set.angle(angle)
	set.offset(width()/2, height()/2)
	set.color(100, 100, 100)
	set.environment()
	
	if distance(mouse.x(), mouse.y(), width()/2, height()/2) < 50px²
		set.color(100, 0, 0)
	end
	circle(0, 0, 50px)
	
	set.color(100, 100, 100)
	
	if distance(mouse.x(), mouse.y(), width()/2, height()/2+75px) < 25px²
		set.color(100, 0, 0)
	end
	rectangle(0, 75px, 50px, 50px)
	set.color(100, 100, 100)
	
	set.angle(-90)
	triangle(-75px, 0, 50px, 50px, 0)
	set.angle(90)
	triangle(75px, 0, 50px, 50px, 0)
	
	image.draw(icon, 0, -75px)
	
	set.decay()
}
