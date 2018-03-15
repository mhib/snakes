import BoardBuilder from './BoardBuilder';

describe('BoardBuilder', () => {
  beforeEach(() => {
    global.document.body.innerHTML = '<div id="parent"></div>';
  });

  afterEach(() => {
    global.document.body.innerHTML = '';
  });

  test('it renders board', () => {
    const width = 30;
    const length = 30;
    BoardBuilder(global.document.getElementById('parent'), width, length);
    expect(global.document.querySelectorAll('.board')).toHaveLength(1);
    expect(global.document.querySelectorAll('.cell')).toHaveLength(width * length);
    expect(global.document.querySelectorAll('.row')).toHaveLength(length);
    expect(global.document.querySelectorAll('.cell[data-x="5"][data-y="5"]')).toHaveLength(1);
  });
});
