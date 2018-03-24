export const connectionUrl = () => (
  window.location.href.replace(/^http/, 'ws')
    .replace('game', 'gamews')
);

export const createSocket = () => new WebSocket(connectionUrl());

export default createSocket;
