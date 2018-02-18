import BoardBuilder from './builders/BoardBuilder';
import EmptyFiller from './fillers/EmptyFiller';
import FruitFiller from './fillers/FruitFiller';
import SnakeFiller from './fillers/SnakeFiller';
import CellSelector from './selectors/CellSelector';

const container = document.getElementById('main-component');

const wsUrl = window.location.href.replace(/(^\w+:|^)\/\//, 'ws://')
  .replace('game', 'gamews');
const socket = new WebSocket(wsUrl);

window.addEventListener('keydown', (event) => {
  event.preventDefault();
  if (event.key === 'ArrowUp') {
    socket.send(JSON.stringify({ direction: 'UP' }));
  } else if (event.key === 'ArrowDown') {
    socket.send(JSON.stringify({ direction: 'DOWN' }));
  } else if (event.key === 'ArrowLeft') {
    socket.send(JSON.stringify({ direction: 'LEFT' }));
  } else if (event.key === 'ArrowRight') {
    socket.send(JSON.stringify({ direction: 'RIGHT' }));
  }
});

let rendered = false;
let selector;

socket.onmessage = (msg) => {
  const data = JSON.parse(msg.data);
  if (!rendered) {
    BoardBuilder(container, data.width, data.length);
    rendered = true;
    selector = new CellSelector(data.width, data.length);
  }
  EmptyFiller(selector.all);
  SnakeFiller(data.snakes, selector.select);
  FruitFiller(data.fruits, selector.select);
};
