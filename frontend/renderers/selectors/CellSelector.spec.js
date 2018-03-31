import BoardBuilder from '../builders/BoardBuilder';
import CellSelector from './CellSelector';

describe('CellSelector', () => {
  const width = 5;
  const length = 5;
  let selector;

  beforeEach(() => {
    global.document.body.innerHTML = '<div id="parent"></div>';
    BoardBuilder(global.document.getElementById('parent'), width, length);
    selector = (new CellSelector(width, length)).select;
  });

  afterEach(() => {
    global.document.body.innerHTML = '';
  });

  describe('existent div', () => {
    test('it returns div', () => {
      const firstDiv = selector(2, 3);
      expect(+firstDiv.dataset.x).toEqual(2);
      expect(+firstDiv.dataset.y).toEqual(3);
      const secondDiv = selector(4, 0);
      expect(+secondDiv.dataset.x).toEqual(4);
      expect(+secondDiv.dataset.y).toEqual(0);
    });
  });

  describe('non-existent div', () => {
    test('it returns undefined', () => {
      expect(selector(-1, 3)).toBeUndefined();
      expect(selector(length + 1, 3)).toBeUndefined();
      expect(selector(3, -1)).toBeUndefined();
      expect(selector(3, length + 1)).toBeUndefined();
    });
  });
});
