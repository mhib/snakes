import BoardRenderer from './DOMBoardRenderer';

describe('DOMBoardRenderer', () => {
  let parent;
  beforeEach(() => {
    global.document.body.innerHTML = '<div id="parent"></div>';
    parent = document.getElementById('parent');
  });

  afterEach(() => {
    global.document.body.innerHTML = '';
  });

  /* eslint-disable no-bitwise */
  const hexToRGB = (hex) => {
    const int = parseInt(hex.slice(1), 16);
    const r = (int >> 16) & 255;
    const g = (int >> 8) & 255;
    const b = int & 255;

    return `rgb(${r}, ${g}, ${b})`;
  };

  describe('#update', () => {
    const firstData = {
      width: 10,
      length: 10,
      snakes: [
        {
          body: [{ x: 5, y: 0 }],
          points: 0,
          name: 'd',
          color: '#B0BC00',
          id: '2bd6aaad-49b2-470b-8881-ad916d94e391',
        },
      ],
      fruits: [],
      state: 1,
    };

    const secondData = {
      width: 10,
      length: 10,
      snakes: [
        {
          body: [{ x: 4, y: 0 }],
          points: 0,
          name: 'd',
          color: '#B0BC00',
          id: '2bd6aaad-49b2-470b-8881-ad916d94e391',
        },
      ],
      fruits: [{ x: 5, y: 7 }],
      state: 1,
    };

    test('it updates DOM', () => {
      const board = new BoardRenderer(parent);
      board.update(firstData);
      expect(parent.getElementsByClassName('row')).toHaveLength(10);
      expect(parent.getElementsByClassName('cell')).toHaveLength(10 * 10);
      expect(parent.querySelector('.cell[data-x="5"][data-y="0"]').style.backgroundColor)
        .toEqual(hexToRGB('#B0BC00'));
      board.update(secondData);
      expect(parent.querySelector('.cell[data-x="4"][data-y="0"]').style.backgroundColor)
        .toEqual(hexToRGB('#B0BC00'));
      expect(parent.querySelector('.cell[data-x="5"][data-y="7"]').style.backgroundColor)
        .toEqual('black');
    });
  });
});
