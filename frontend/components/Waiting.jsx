import React from 'react';
import styled from 'styled-components';

const Paragraph = styled.p`
font-weight: bold;
`;

const handleFocus = event => event.target.select();

const Waiting = () => (
  <div>
    <Paragraph>Waiting for other players</Paragraph>
    <span>Link to join for other players: </span>
    <input type="text" onFocus={handleFocus} defaultValue={window.location.href} />
  </div>
);

export default Waiting;
