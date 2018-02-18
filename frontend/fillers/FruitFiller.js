const FruitFiller = (fruits, selector) => {
  fruits.forEach((fruit) => {
    const cell = selector(fruit.x, fruit.y);
    cell.style.backgroundColor = 'black';
  });
};

export default FruitFiller;
