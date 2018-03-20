import ColorHash from 'color-hash';
import memoize from 'lodash.memoize';
import FillCell from './FillCell';

const colorHash = new ColorHash();
const memoizedHex = memoize(colorHash.hex.bind(colorHash));

const SnakeFiller = (snakes, selector) => {
  snakes.forEach(({ add, remove, id }) => {
    if (add) {
      FillCell(add, memoizedHex(id), selector);
    }
    if (remove) {
      FillCell(remove, 'initial', selector);
    }
  });
};

export default SnakeFiller;
