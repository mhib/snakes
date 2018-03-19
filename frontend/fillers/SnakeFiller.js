import ColorHash from 'color-hash';
import memoize from 'lodash.memoize';

const colorHash = new ColorHash();
const memoizedHex = memoize(colorHash.hex.bind(colorHash));

const fillCell = (point, color, selector) => {
  const cell = selector(point.x, point.y);
  if (cell) {
    cell.style.backgroundColor = color;
  }
};

const SnakeFiller = (snakes, selector) => {
  snakes.forEach(({ add, remove, id }) => {
    const hex = memoizedHex(id);
    if (add) {
      fillCell(add, hex, selector);
    }
    if (remove) {
      fillCell(remove, 'initial', selector);
    }
  });
};

export default SnakeFiller;
