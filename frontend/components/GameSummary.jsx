import React from 'react';
import styled from 'styled-components';
import GameSummaryShape from './shapes/GameSummary';

const Container = styled.div`
display: block;
border-radius: 3px;
border: 1px solid #DDDDDD;
cursor: pointer;
padding: 5px;
&:hover {
  background-color: #EEEEEE;
}
`;

export default class GameSummary extends React.Component {
  static propTypes = {
    game: GameSummaryShape.isRequired,
  };

  constructor(props) {
    super(props);
    this.handleClick = this.handleClick.bind(this);
  }

  gameHref() {
    return `${window.location.protocol}//${window.location.host}/game/${this.props.game.id}`;
  }

  handleClick() {
    window.location.href = this.gameHref();
  }

  botsParagraph() {
    if (this.props.game.bots > 0) {
      return <p>Bots: {this.props.game.bots}</p>;
    }
    return null;
  }

  render() {
    const {
      connected, players, width, length, moveTick, foodTick,
    } = this.props.game;
    return (
      <Container onClick={this.handleClick}>
        <p>Players: {connected}/{players}</p>
        {this.botsParagraph()}
        <p>Size: {width}x{length}</p>
        <p>Tick(move/food): {moveTick}/{foodTick} ms</p>
      </Container>
    );
  }
}
