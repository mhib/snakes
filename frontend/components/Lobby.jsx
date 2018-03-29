import React from 'react';
import styled from 'styled-components';
import LobbyForm from './LobbyForm';

const Container = styled.div`
display: grid;
padding: 20px;
grid-template-columns: 4fr 8fr;
grid-template-rows: 100%;
`;

const Lobby = () => (
  <Container>
    <LobbyForm />
    <div />
  </Container>
);

export default Lobby;
