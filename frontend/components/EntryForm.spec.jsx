import React from 'react';
import { render } from 'enzyme';
import EntryForm from './EntryForm';

describe('<EntryForm />', () => {
  let wrapper;
  const props = {
    onSubmit() {
    },
  };
  beforeEach(() => {
    wrapper = render(<EntryForm {...props} />);
  });

  describe('#render', () => {
    test('it renders form', () => {
      expect(wrapper.find('label')).toHaveLength(2);
    });
  });
});
