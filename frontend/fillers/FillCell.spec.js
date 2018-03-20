import FillCell from './FillCell';

describe('BoardBuilder', () => {
  let div;
  beforeEach(() => {
    global.document.body.innerHTML = '<div id="parent"></div>';
    div = document.getElementById('parent');
  });

  afterEach(() => {
    global.document.body.innerHTML = '';
  });

  describe('point', () => {
    test('it fills cell', () => {
      const selector = () => div;
      const color = 'rgb(123, 127, 89)';
      FillCell({ x: 1, y: 1 }, color, selector);
      expect(div.style.backgroundColor).toEqual(color);
    });
  });

  describe('no point', () => {
    test('it fills cell', () => {
      const selector = () => null;
      const color = 'rgb(0, 0, 0)';
      expect(() => {
        FillCell({ x: 1, y: 1 }, color, selector);
      }).not.toThrowError();
    });
  });
});
