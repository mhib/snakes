import React from 'react';
import styled from 'styled-components';
import Cell from './Cell';

const RowDiv = styled.div`
flex: 1 1 auto
display: flex;
flex-direction: row;
height: 10px;
`;

const Row = ({board, y}) => (
  <RowDiv>
  {[...Array(board.width)].map((_, i) => (
    <Cell board={board} y={y} x={i} key={i} />
  ))}
  </RowDiv>
);

export default Row;
