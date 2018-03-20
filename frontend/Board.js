import BoardBuilder from './builders/BoardBuilder';
import FruitFiller from './fillers/FruitFiller';
import SnakeFiller from './fillers/SnakeFiller';
import CellSelector from './selectors/CellSelector';
import PointEquality from './equalities/PointEquality';

export default class Board {
  constructor(container) {
    this.container = container;
    this.rendered = false;
    this.updateSnake = this.updateSnake.bind(this);
  }

  update(data) {
    if (!this.rendered) {
      this.render(data);
    }
    SnakeFiller(this.generateSnakeDiff(data), this.selector.select);
    FruitFiller(data.fruits, this.selector.select);
  }

  render(data) {
    BoardBuilder(this.container, data.width, data.length);
    this.rendered = true;
    this.selector = new CellSelector(data.width, data.length);
    this.snakes = {};
    SnakeFiller(data.snakes.map(({ body, id }) => ({ id, add: body[0] })), this.selector.select);
    FruitFiller(data.fruits, this.selector.select);
    data.snakes.forEach(this.updateSnake);
  }

  generateSnakeDiff(data) {
    return data.snakes.map((snake) => {
      const update = { id: snake.id };
      const oldSnake = this.snakes[snake.id];
      if (!PointEquality(oldSnake.body[0], snake.body[0])) {
        [update.add] = snake.body;
      }
      if (
        !PointEquality(
          oldSnake.body[oldSnake.body.length - 1],
          snake.body[snake.body.length - 1],
        )
      ) {
        update.remove = oldSnake.body[oldSnake.body.length - 1];
      }
      this.updateSnake(snake);
      return update;
    });
  }

  updateSnake(snake) {
    this.snakes[snake.id] = snake;
  }
}
