const FillCell = (point, color, selector) => {
  const cell = selector(point.x, point.y);
  if (cell) {
    cell.style.backgroundColor = color;
  }
};

export default FillCell;
