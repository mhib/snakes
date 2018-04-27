import PropTypes from 'prop-types';

const GameSummary = PropTypes.shape({
  id: PropTypes.string.isRequired,
  players: PropTypes.number.isRequired,
  connected: PropTypes.number.isRequired,
  width: PropTypes.number.isRequired,
  length: PropTypes.number.isRequired,
  foodTick: PropTypes.number.isRequired,
  moveTick: PropTypes.number.isRequired,
  bots: PropTypes.number.isRequired,
});

export default GameSummary;
