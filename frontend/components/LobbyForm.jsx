import React from 'react';

export default class LobbyForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      players: 1,
      move_tick: 100,
      food_tick: 10000,
      width: 10,
      length: 10,
      endOnLastPlayer: true,
    };
    this.handleInputChange = this.handleInputChange.bind(this);
  }

  handleInputChange({ target }) {
    const {
      name, type, checked, valueAsNumber, value,
    } = target;
    let inputValue;
    switch (type) {
      case 'number':
        inputValue = valueAsNumber || value;
        break;
      case 'checkbox':
        inputValue = checked;
        break;
      default:
        inputValue = value;
    }

    this.setState({
      [name]: inputValue,
    });
  }

  lastPlayerInput() {
    if (this.state.players <= 1) {
      return null;
    }
    return (
      <div>
        End game with only 1 player:
        <input
          type="checkbox"
          checked={this.state.endOnLastPlayer}
          onChange={this.handleInputChange}
        />
      </div>
    );
  }

  render() {
    return (
      <div>
        <h3>Create new game</h3>
        <form action="/new_game/" method="POST">
          Number of players:
          <input value={this.state.players} onChange={this.handleInputChange} name="players" required type="number" min="1" max="30" /><br />
          Move tick in milliseconds:
          <input value={this.state.move_tick} onChange={this.handleInputChange} name="move_tick" required type="number" min="1" max="3000" /><br />
          Food tick in milliseconds:
          <input value={this.state.food_tick} onChange={this.handleInputChange} name="food_tick" required type="number" min="0" max="120000" /><br />
          Width:
          <input value={this.state.width} onChange={this.handleInputChange} name="width" required type="number" min="10" max="100" /><br />
          Length:
          <input value={this.state.length} onChange={this.handleInputChange} name="length" required type="number" min="10" max="100" /><br />
          {this.lastPlayerInput()}
          <input type="submit" value="New game" readOnly />
        </form>
      </div>
    );
  }
}
