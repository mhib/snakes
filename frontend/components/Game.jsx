import React from 'react';
import styled from 'styled-components';
import bindAll from 'lodash.bindall';
import EntryForm from './EntryForm';
import Waiting from './Waiting';
import Ranking from './Ranking';
import SocketFactory from '../factories/SocketFactory';
import Board from '../Board';

const GameContainer = styled.div`
position: absolute;
width: 100%;
height: 100%;
display: flex;
justify-content: center;
padding: 2em;
`;

const BoardContainer = styled.div`
display: flex;
flex-direction: column;
`;

export default class Game extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      gameState: 'notConnected',
      ranking: [],
    };
    bindAll(this, ['handleSubmit', 'handleUpdate', 'handleClose', 'handleKeyDown',
      'updateRanking']);
  }

  shouldRenderEntryForm() {
    return this.state.gameState === 'notConnected';
  }

  shouldRenderWaiting() {
    return this.state.gameState === 'waiting';
  }

  shouldRenderBoard() {
    return this.state.gameState === 'playing';
  }

  shouldRenderRanking() {
    return this.state.ranking.length !== 0;
  }

  prepareSocket(formState) {
    this.socket = SocketFactory();
    this.socket.onopen = () => {
      this.socket.send(JSON.stringify(formState));
    };
    this.socket.onmessage = this.handleUpdate;
    this.socket.onclose = this.handleClose;
  }

  handleSubmit(formState) {
    this.board = new Board(this.boardDiv);
    this.prepareSocket(formState);
    window.addEventListener('keydown', this.handleKeyDown);
    this.setState({ gameState: 'waiting' });
  }

  handleClose() {
    this.setState({ gameState: 'ended' });
  }

  handleUpdate({ data }) {
    if (this.state.gameState !== 'playing') {
      this.setState({ gameState: 'playing' });
    }
    const parsedData = JSON.parse(data);
    this.board.update(parsedData);
    this.updateRanking(parsedData);
  }

  handleKeyDown(event) {
    event.preventDefault();
    if (event.key === 'ArrowUp') {
      this.socket.send(JSON.stringify({ direction: 'UP' }));
    } else if (event.key === 'ArrowDown') {
      this.socket.send(JSON.stringify({ direction: 'DOWN' }));
    } else if (event.key === 'ArrowLeft') {
      this.socket.send(JSON.stringify({ direction: 'LEFT' }));
    } else if (event.key === 'ArrowRight') {
      this.socket.send(JSON.stringify({ direction: 'RIGHT' }));
    }
  }

  updateRanking(data) {
    const snakes = data.snakes.concat().sort((l, r) => r.points - l.points);
    this.setState({ ranking: snakes });
  }

  render() {
    return (
      <GameContainer>
        {this.shouldRenderEntryForm() && <EntryForm onSubmit={this.handleSubmit} />}
        {this.shouldRenderWaiting() && <Waiting />}
        {this.shouldRenderRanking() && <Ranking snakes={this.state.ranking} />}
        <BoardContainer innerRef={(div) => { this.boardDiv = div; }} />
      </GameContainer>
    );
  }
}