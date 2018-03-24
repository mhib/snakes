import FillCell from './FillCell';

const SnakeFiller = (snakes, selector) => {
  snakes.forEach(({ add, remove, color }) => {
    if (add) {
      FillCell(add, color, selector);
    }
    if (remove) {
      FillCell(remove, 'initial', selector);
    }
  });
};

export default SnakeFiller;
