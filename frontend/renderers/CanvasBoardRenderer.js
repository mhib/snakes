import PointEquality from '../equalities/PointEquality';

const CELL_WIDTH = 14;
const LINE_WIDTH = 2;

export default class Board {
  constructor(domElement, initialBoard) {
    this.initialBoard = initialBoard;
    this.domElement = domElement;
    this.ctx = domElement.getContext('2d');
    this.height = (initialBoard.length * CELL_WIDTH) + LINE_WIDTH;
    this.width = (initialBoard.width * CELL_WIDTH) + LINE_WIDTH;
    this.snakes = {};
    this.rendered = false;
  }

  setDomDimensions() {
    this.domElement.style.height = this.height;
    this.domElement.height = this.height;
    this.domElement.style.width = this.height;
    this.domElement.width = this.width;
  }

  update(data) {
    if (!this.rendered) {
      this.render(data);
      this.rendered = true;
      return;
    }
    this.updateSnakes(data);
    this.updateFruits(data);
  }

  static getRect(x, y) {
    return [
      LINE_WIDTH + (x * CELL_WIDTH),
      LINE_WIDTH + (y * CELL_WIDTH),
      CELL_WIDTH - LINE_WIDTH,
      CELL_WIDTH - LINE_WIDTH,
    ];
  }

  fillRect(x, y, color) {
    this.ctx.fillStyle = color;
    this.ctx.fillRect(...this.constructor.getRect(x, y));
    this.ctx.stroke();
  }

  updateSnakes(data) {
    data.snakes.forEach((snake) => {
      const oldSnake = this.snakes[snake.id];
      if (!PointEquality(oldSnake.body[0], snake.body[0])) {
        this.fillRect(snake.body[0].x, snake.body[0].y, snake.color);
      }
      if (
        !PointEquality(
          oldSnake.body[oldSnake.body.length - 1],
          snake.body[snake.body.length - 1],
        )
      ) {
        this.fillRect(oldSnake.body[oldSnake.body.length - 1].x, oldSnake.body[oldSnake.body.length - 1].y, 'white');
      }
      this.updateSnake(snake);
    });
  }

  updateFruits({ fruits }) {
    fruits.forEach(({ x, y }) => {
      this.fillRect(x, y, 'black');
    });
  }

  render(data) {
    this.drawBoard();
    this.updateFruits(data);
    data.snakes.forEach((snake) => {
      this.fillRect(snake.body[0].x, snake.body[0].y, snake.color);
      this.updateSnake(snake);
    });
  }

  drawBoard() {
    this.setDomDimensions();

    this.ctx.fillStyle = 'black';
    this.ctx.fillRect(0, 0, this.width, this.height);
    this.ctx.stroke();

    this.ctx.fillStyle = 'white';
    this.ctx.lineWidth = LINE_WIDTH / 2;
    for (let x = 0; x < this.initialBoard.width; x += 1) {
      for (let y = 0; y < this.initialBoard.length; y += 1) {
        this.ctx.fillRect(...this.constructor.getRect(x, y));
      }
    }
    this.ctx.stroke();
  }

  updateSnake(snake) {
    this.snakes[snake.id] = snake;
  }
}
