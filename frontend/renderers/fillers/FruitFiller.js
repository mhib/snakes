import FillCell from './FillCell';

const FruitFiller = (fruits, selector) => {
  fruits.forEach((fruit) => {
    FillCell(fruit, 'black', selector);
  });
};

export default FruitFiller;
