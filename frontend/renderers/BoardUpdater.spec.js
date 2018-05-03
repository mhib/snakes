import BoardUpdater from './BoardUpdater';
import CanvasRenderer from './CanvasBoardRenderer';

jest.mock('./CanvasBoardRenderer');

describe('BoardUpdater', () => {
  beforeEach(() => {
    CanvasRenderer.mockClear();
  });

  describe('#update', () => {
    const initialData = JSON.parse(`{
      "width":10,"length":10,"snakes":[{"body":[{"x":5,"y":0}],"points":0,"name":"d","color":"#B0BC00","id":"2bd6aaad-49b2-470b-8881-ad916d94e391"}],"fruits":[],"state":1
    }`);

    const secondData = JSON.parse(`{
      "width":10,"length":10,"snakes":[{"body":[{"x":4,"y":0}],"points":0,"name":"d","color":"#B0BC00","id":"2bd6aaad-49b2-470b-8881-ad916d94e391"}],"fruits":[{"x":5,"y":7}],"state":2
    }`);

    it('renders initial board on first update', () => {
      const renderer = new CanvasRenderer();
      const updater = new BoardUpdater(renderer);

      updater.update(initialData);
      expect(renderer.drawBoard).toHaveBeenCalledTimes(1);
      expect(renderer.fillRect).toHaveBeenCalledWith(5, 0, '#B0BC00');
    });

    it('updates correctly on snake move', () => {
      const renderer = new CanvasRenderer();
      const updater = new BoardUpdater(renderer);

      updater.update(initialData);
      renderer.fillRect.mockClear();
      updater.update(secondData);

      expect(renderer.fillRect).toHaveBeenCalledWith(5, 0, 'white');
      expect(renderer.fillRect).toHaveBeenCalledWith(4, 0, '#B0BC00');
      expect(renderer.fillRect).toHaveBeenCalledWith(4, 0, '#B0BC00');
      expect(renderer.fillRect).toHaveBeenCalledWith(5, 7, 'black');
    });
  });
});
