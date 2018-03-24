import React from 'react';
import styled from 'styled-components';
import PropTypes from 'prop-types';
import SnakeShape from './shapes/Snake';

const RankingContainer = styled.div`
position: fixed;
top: 30px;
left: 30px;
background-color: #DDDDDD;
opacity: 0.7;
padding: 1em;
border-radius: 5px;
&:hover {
  opacity: 1;
}
`;

const UserName = styled.span`
text-shadow: 0px 0px 1px black;
`;

const Ranking = ({ snakes }) => (
  <RankingContainer>
    {
      snakes.map((snake, idx) => (
        <div key={snake.id}>
          {idx + 1}. <UserName style={{ color: snake.color }}>{snake.name}</UserName>
          - {snake.points} points
        </div>
      ))
    }
  </RankingContainer>
);

Ranking.propTypes = {
  snakes: PropTypes.arrayOf(SnakeShape).isRequired,
};

export default Ranking;
