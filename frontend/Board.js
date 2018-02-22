import BoardBuilder from './builders/BoardBuilder';
import EmptyFiller from './fillers/EmptyFiller';
import FruitFiller from './fillers/FruitFiller';
import SnakeFiller from './fillers/SnakeFiller';
import CellSelector from './selectors/CellSelector';

export default class Board {
  constructor(container) {
    this.container = container;
    this.rendered = false;
  }

  update(data) {
    if (!this.rendered) {
      this.render(data);
    }
    EmptyFiller(this.selector.all);
    SnakeFiller(data.snakes, this.selector.select);
    FruitFiller(data.fruits, this.selector.select);
  }

  render(data) {
    BoardBuilder(this.container, data.width, data.length);
    this.rendered = true;
    this.selector = new CellSelector(data.width, data.length);
  }
}
