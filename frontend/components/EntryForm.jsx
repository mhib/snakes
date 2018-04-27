import React from 'react';
import styled from 'styled-components';
import CompactPicker from 'react-color/lib/Compact';
import PropTypes from 'prop-types';

const FormTag = styled.form`
display: flex;
flex-direction: column;
`;

const LabelTag = styled.label`
font-weight: bold;
padding-top: 10px;
`;

const SubmitTag = styled.input`
margin-top: 10px;
`;

const defaultColors = [
  '#F44E3B', '#FE9200', '#FCDC00', '#DBDF00', '#A4DD00', '#68CCCA', '#73D8FF', '#AEA1FF',
  '#FDA1FF', '#D33115', '#E27300', '#FCC400', '#B0BC00', '#68BC00', '#16A5A5',
  '#009CE0', '#7B64FF', '#FA28FF', '#9F0500', '#C45100', '#FB9E00', '#808900', '#194D33',
  '#0C797D', '#0062B1', '#653294', '#AB149E',
];

const shuffleArray = arr => (
  arr
    .map(a => [Math.random(), a])
    .sort((a, b) => a[0] - b[0])
    .map(a => a[1])
);

export default class EntryForm extends React.Component {
  static propTypes = EntryForm.propTypes = {
    onSubmit: PropTypes.func.isRequired,
  };

  constructor(props) {
    super(props);
    this.colors = shuffleArray(defaultColors);
    this.state = {
      name: '',
      color: this.colors[0],
    };
    this.handleNameChange = this.handleNameChange.bind(this);
    this.handleColorChange = this.handleColorChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleNameChange(event) {
    this.setState({ name: event.target.value });
  }

  handleColorChange({ hex }) {
    if (hex !== '#000000' && hex !== '#FFFFFF') {
      this.setState({ color: hex });
    }
  }

  handleSubmit(event) {
    event.preventDefault();
    this.props.onSubmit(this.state);
  }

  render() {
    return (
      <FormTag onSubmit={this.handleSubmit}>
        <LabelTag htmlFor="name-input">
        Name:
        </LabelTag>
        <input id="name-input" type="text" value={this.state.name} onChange={this.handleNameChange} />
        <LabelTag htmlFor="color-input">
        Color:
        </LabelTag>
        <CompactPicker id="color-input" colors={this.colors} color={this.state.color} onChange={this.handleColorChange} />
        <SubmitTag type="submit" value="Join game" />
      </FormTag>
    );
  }
}

