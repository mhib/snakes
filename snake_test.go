package main

import "testing"

func TestChangeDirection(t *testing.T) {
	snake := Snake{[]Point{}, 0, LEFT, false, "dd"}
	snake.changeDirection(RIGHT)
	if snake.Direction == RIGHT {
		t.Error("Changed to opposite direction")
	}
	snake.changeDirection(UP)
	if snake.Direction != UP {
		t.Error("Cannot change direction")
	}
}

func TestGrow(t *testing.T) {
	snake := Snake{[]Point{}, 0, LEFT, false, "dd"}
	snake.grow(3)
	if len(snake.Body) != 3 {
		t.Error("Growing failed")
	}
}

func TestIncludes(t *testing.T) {
	snake := Snake{[]Point{Point{1, 1}, Point{1, 2}}, 0, LEFT, false, "dd"}
	if !snake.includes(Point{1, 1}) {
		t.Error("Includes failed")
	}
	if snake.includes(Point{1, 3}) {
		t.Error("Includes failed")
	}
}

func TestCollideItself(t *testing.T) {
	first := Snake{[]Point{Point{1, 1}, Point{1, 2}, Point{1, 1}}, 0, LEFT, false, "dd"}
	if !first.collideItself() {
		t.Error("Includes failed")
	}
	second := Snake{[]Point{Point{1, 1}, Point{1, 2}}, 0, LEFT, false, "dd"}
	if second.collideItself() {
		t.Error("Includes failed")
	}
}

func TestMove(t *testing.T) {
	first := Snake{[]Point{Point{0, 2}, Point{1, 2}, Point{2, 2}}, 0, LEFT, false, "dd"}
	first.move(5, 5)
	newHead := Point{4, 2}
	if first.head() != newHead {
		t.Error("Move failed")
	}
}
