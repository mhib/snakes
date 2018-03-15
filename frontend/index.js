import Board from './Board';

const board = new Board(document.getElementById('main-component'));

const wsUrl = window.location.href.replace(/^http/, 'ws')
  .replace('game', 'gamews');
const socket = new WebSocket(wsUrl);
socket.onmessage = (msg) => {
  board.update(JSON.parse(msg.data));
};

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

