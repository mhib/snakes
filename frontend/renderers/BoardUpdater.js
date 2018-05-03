import PointEquality from '../equalities/PointEquality';

export default class BoardUpdater {
  constructor(renderer) {
    this.snakes = {};
    this.wasRendered = false;
    this.renderer = renderer;
  }

  update(data) {
    if (!this.wasRendered) {
      this.render(data);
      this.wasRendered = true;
      return;
    }
    this.updateSnakes(data);
    this.updateFruits(data);
  }

  render(data) {
    this.renderer.drawBoard();
    this.updateFruits(data);
    data.snakes.forEach((snake) => {
      this.renderer.fillRect(snake.body[0].x, snake.body[0].y, snake.color);
      this.updateSnake(snake);
    });
  }

  updateSnakes(data) {
    data.snakes.forEach((snake) => {
      const oldSnake = this.snakes[snake.id];
      if (!PointEquality(oldSnake.body[0], snake.body[0])) {
        this.renderer.fillRect(snake.body[0].x, snake.body[0].y, snake.color);
      }
      if (
        !PointEquality(
          oldSnake.body[oldSnake.body.length - 1],
          snake.body[snake.body.length - 1],
        )
      ) {
        this.renderer.fillRect(oldSnake.body[oldSnake.body.length - 1].x, oldSnake.body[oldSnake.body.length - 1].y, 'white');
      }
      this.updateSnake(snake);
    });
  }

  updateFruits({ fruits }) {
    fruits.forEach(({ x, y }) => {
      this.renderer.fillRect(x, y, 'black');
    });
  }

  updateSnake(snake) {
    this.snakes[snake.id] = snake;
  }
}
