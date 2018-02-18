import ColorHash from 'color-hash';
import memoize from 'lodash.memoize';

const colorHash = new ColorHash();
const memoizedHex = memoize(colorHash.hex.bind(colorHash));

const SnakeFiller = (snakes, selector) => {
  snakes.forEach((snake) => {
    const hex = memoizedHex(snake.id);
    snake.body.forEach((point) => {
      const cell = selector(point.x, point.y);
      if (cell) {
        cell.style.backgroundColor = hex;
      }
    });
  });
};

export default SnakeFiller;
