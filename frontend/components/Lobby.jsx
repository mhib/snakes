import React from 'react';
import styled from 'styled-components';
import LobbyForm from './LobbyForm';
import GameList from './GameList';
import SocketFactory from '../factories/LobbySocketFactory';

const Container = styled.div`
display: grid;
padding: 20px;
grid-template-columns: 4fr 8fr;
grid-template-rows: 100%;
`;

export default class Lobby extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      games: [],
    };
    this.prepareSocket();
    this.closeSocket = this.closeSocket.bind(this);
  }

  prepareSocket() {
    this.socket = SocketFactory();
    this.socket.onmessage = ({ data }) => {
      this.setState({ games: JSON.parse(data) });
    };
  }

  closeSocket() {
    this.socket.close();
  }

  render() {
    return (
      <Container>
        <LobbyForm onSubmit={this.closeSocket} />
        <GameList games={this.state.games} />
      </Container>
    );
  }
}
