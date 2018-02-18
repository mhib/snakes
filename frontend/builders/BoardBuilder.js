const newCell = (x, y) => {
  const square = document.createElement('div');
  square.setAttribute('class', 'cell');
  square.setAttribute('data-x', x);
  square.setAttribute('data-y', y);
  return square;
};

const BoardBuilder = (parent, width, length) => {
  const board = document.createElement('div');
  board.setAttribute('class', 'board');
  for (let y = 0; y < length; y += 1) {
    const row = document.createElement('div');
    row.setAttribute('class', 'row');
    for (let x = 0; x < width; x += 1) {
      row.appendChild(newCell(x, y));
    }
    parent.appendChild(row);
  }
  parent.appendChild(board);
};

export default BoardBuilder;
