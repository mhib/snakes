import React from 'react';
import { arrayOf } from 'prop-types';
import styled from 'styled-components';
import GameSummaryShape from './shapes/GameSummary';
import GameSummary from './GameSummary';

const Container = styled.div`
display: flex;
flex-direction: column;
align-items: stretch;
text-align: center;
`;

const gameListBody = (games) => {
  if (games.length > 0) {
    return games.map(game => (
      <GameSummary key={game.id} game={game} />
    ));
  }
  return <p>No open games.</p>;
};

const GameList = ({ games }) => (
  <Container>
    <h3>Open games:</h3>
    {gameListBody(games)}
  </Container>
);

GameList.propTypes = {
  games: arrayOf(GameSummaryShape).isRequired,
};

export default GameList;
