import React from 'react';
import { render, mount } from 'enzyme';
import EntryForm from './EntryForm';

describe('<EntryForm />', () => {
  const props = {
    onSubmit: jest.fn(),
  };

  beforeEach(() => {
    props.onSubmit.mockClear();
  });


  test('it renders form', () => {
    const wrapper = render(<EntryForm {...props} />);
    expect(wrapper.find('label')).toHaveLength(2);
  });

  describe('form submit', () => {
    const userName = 'Name';
    let wrapper;

    beforeEach(() => {
      wrapper = mount(<EntryForm {...props} />);
      wrapper.setState({ name: userName });
      wrapper.find('form').simulate('submit');
    });

    it('calls onSubmit with wrapper state', () => {
      expect(props.onSubmit.mock.calls.length).toBe(1);
      expect(props.onSubmit.mock.calls[0][0]).toEqual(wrapper.state());
    });
  });
});
