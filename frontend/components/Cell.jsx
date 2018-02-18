import React from 'react';
import styled from 'styled-components';
import PointSerializer from './PointSerializer';

const CellDiv = styled.div`
border: 1px solid black;
width: 10px;
height: 10px;
display: inline-block;
`;

export default class Cell extends React.Component {
  constructor(props) {
    super(props);
    this.pointID = PointSerializer({x: props.x, y: props.y})
  }

  render() {
    return (<CellDiv style={{backgroundColor: this._generateBackgroundColor()}} />)
  }

  _generateBackgroundColor() {
    let color;
    if (this.props.board.fruits.has(this.pointID)) {
      return 'black';
    } else if (color = this.props.board.snakeCells.get(this.pointID))  {
      return color;
    } else {
      return 'initial';
    }
  }
}
