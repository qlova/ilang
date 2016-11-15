type Car {
    speed
}

method new(Car) {
    speed = 100
}

method drive(Car) {
    print("My top speed is: ", speed)
}

method stop(Car) {
	speed = 0
}

type SportsCar is Car {}

method new(SportsCar) {
    speed = 200
}

method drive(SportsCar) {
    print("Vroom vroom!")
    drive(Car)
}

software {
    var car = new SportsCar()
    drive(car) //--> Should output "Vroom vroom! My top speed is: 200"
    stop(car)
    drive(car)
}
