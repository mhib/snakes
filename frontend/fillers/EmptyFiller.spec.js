import EmptyFiller from './EmptyFiller';

describe('EmptyFiller', () => {
  const DIV_COUNT = 100;
  beforeEach(() => {
    for (let counter = 0; counter < DIV_COUNT; counter += 1) {
      const div = global.document.createElement('div');
      const stub = jest.fn();
      div.setAttribute('class', 'match');
      Object.defineProperty(div.style, 'backgroundColor', {
        set: e => stub(e),
      });
      div.style.stub = stub;
      global.document.body.appendChild(div);
    }
  });

  afterEach(() => {
    global.document.body.innerHTML = '';
    global.document.head.innerHTML = '';
  });

  test('it renders board', () => {
    expect([...global.document.querySelectorAll('.match')]
      .filter(el => el.style.backgroundColor !== 'initial')).toHaveLength(DIV_COUNT);
    EmptyFiller([...global.document.querySelectorAll('.match')]);
    expect([...global.document.querySelectorAll('.match')]
      .every(el => el.style.stub.mock.calls[0][0] === 'initial')).toBe(true);
  });
});
