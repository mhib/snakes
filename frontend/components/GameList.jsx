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

const GameList = ({ games }) => (
  <Container>
    <h3>Open games:</h3>
    {
      games.map(game => (
        <GameSummary key={game.id} game={game} />
      ))
    }
  </Container>
);

GameList.propTypes = {
  games: arrayOf(GameSummaryShape).isRequired,
};

export default GameList;
