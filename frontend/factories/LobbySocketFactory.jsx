export const connectionUrl = () => {
  const protocol = window.location.protocol.replace(/^http/, 'ws');
  return `${protocol}//${window.location.host}/lobby`;
};

export const createSocket = () => new WebSocket(connectionUrl());

export default createSocket;
