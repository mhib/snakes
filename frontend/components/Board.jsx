import React from 'react';
import styled from 'styled-components';
import ColorHash from 'color-hash';
import memoize from 'lodash.memoize';
import Row from './Row';
import PointSerializer from './PointSerializer'

const Wrapper = styled.div`
display: flex;
flex-direction: column;
`;

const colorHash = new ColorHash();
const memoizedHex = memoize(colorHash.hex.bind(colorHash));

export default class Board extends React.Component {
  constructor(props) {
    super(props)
    this.socket = new WebSocket(this._websocketUrl())
    this.socket.onmessage = this._handleWebsocket.bind(this)
    this.state = {
      width: 0,
      length: 0,
      snakeCells: new Map([]),
      fruits: new Set([]),
      boardState: 0
    }
    this._bindKeys();
  }

  render() {
    return (
      <Wrapper>
      {this._isWaiting() && <h1>Waiting for other players</h1>}
      {[...Array(this.state.length)].map((_, i) => (
        <Row board={this.state} key={i} y={i} />
      ))}
      {this._isPreparing() && <h1>Prepare yourself</h1>}
      </Wrapper>
    )
  }

  _isWaiting() {
    return this.state.boardState === 0;
  }

  _isPreparing() {
    return this.state.boardState === 1;
  }

  _handleWebsocket(message) {
    this.setState(() => {
      const data = JSON.parse(message.data);
      const { width, length } = data;
      const boardState = data.state;
      const fruits = new Set(data.fruits.map(PointSerializer));
      const snakeCells = new Map([].concat(...data.snakes.map((s) => {
        const hex = memoizedHex(s.id);
        return s.body.map((point) => [PointSerializer(point), hex]);
      })))
      return { width, length, boardState, fruits, snakeCells };
    })
  }

  _websocketUrl() {
    return window.location.href.replace(/(^\w+:|^)\/\//, 'ws://')
      .replace("game", "gamews")
  }

  _bindKeys() {
    window.addEventListener('keydown', (event) => {
      event.preventDefault();
      if (event.key === "ArrowUp") {
        this.socket.send(JSON.stringify({ direction: "UP" }));
      } else if (event.key === "ArrowDown") {
        this.socket.send(JSON.stringify({ direction: "DOWN" }));
      } else if (event.key === "ArrowLeft") {
        this.socket.send(JSON.stringify({ direction: "LEFT" }));
      } else if (event.key === "ArrowRight") {
        this.socket.send(JSON.stringify({ direction: "RIGHT" }));
      }
    })
  }
}
