import React from 'react';
import { mount, shallow } from 'enzyme';
import Game from './Game';
import EntryForm from './EntryForm';
import Ranking from './Ranking';
import Board from '../renderers/BoardUpdater';

const mockSocket = { send: jest.fn() };
jest.mock('../factories/GameSocketFactory', () => () => mockSocket);
jest.mock('../renderers/BoardUpdater');

describe('<Game />', () => {
  beforeEach(() => {
    mockSocket.send.mockClear();
    Board.mockClear();
  });
  describe('#render', () => {
    describe('initial state', () => {
      let wrapper;
      beforeEach(() => {
        wrapper = shallow(<Game />);
      });
      test('it renders correctly', () => {
        expect(wrapper.find(EntryForm)).toHaveLength(1);
        expect(wrapper.find(Ranking)).toHaveLength(0);
      });
    });
  });

  describe('game scenario', () => {
    let wrapper;
    const initialMessage = '{"width":10,"length":10,"snakes":[{"body":[{"x":5,"y":0}],"points":0,"name":"d","color":"#B0BC00","id":"2bd6aaad-49b2-470b-8881-ad916d94e391"}],"fruits":[],"state":1}';
    const secondMessage = '{"width":10,"length":10,"snakes":[{"body":[{"x":4,"y":0},{"x":5,"y":0}],"points":0,"name":"d","color":"#B0BC00","id":"2bd6aaad-49b2-470b-8881-ad916d94e391"}],"fruits":[{"x":5,"y":7}],"state":2}';
    beforeEach(() => {
      wrapper = mount(<Game />);
    });
    test('rerenders', () => {
      // form submits
      const formState = { name: 'Dd', color: '#123456' };
      wrapper.instance().handleSubmit(formState);
      expect(mockSocket.onmessage).toBeDefined();
      expect(mockSocket.onopen).toBeDefined();
      mockSocket.send.mockClear();
      mockSocket.onopen();
      wrapper.update();
      expect(mockSocket.send.mock.calls[0][0]).toEqual(JSON.stringify(formState));
      expect(wrapper.find(EntryForm)).toHaveLength(0);

      // First update
      mockSocket.onmessage({ data: initialMessage });
      wrapper.update();
      const boardInstance = Board.mock.instances[0];
      expect(boardInstance.update.mock.calls[0][0]).toEqual(JSON.parse(initialMessage));
      expect(wrapper.find(Ranking)).toHaveLength(1);
      expect(wrapper.find(Ranking).prop('snakes'))
        .toEqual(wrapper.instance().state.ranking);

      // Second update
      mockSocket.onmessage({ data: secondMessage });
      wrapper.update();
      expect(boardInstance.update.mock.calls[1][0]).toEqual(JSON.parse(secondMessage));

      // Messages
      mockSocket.send.mockClear();
      ['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight']
        .forEach(key => window.dispatchEvent(new window.KeyboardEvent('keydown', { key })));
      ['UP', 'DOWN', 'LEFT', 'RIGHT']
        .forEach((direction, idx) => {
          expect(mockSocket.send.mock.calls[idx][0])
            .toEqual(JSON.stringify({ direction }));
        });

      // close
      mockSocket.onclose();
      wrapper.update();
    });
  });
});
