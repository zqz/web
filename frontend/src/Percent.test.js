import React from 'react';
import {shallow} from 'enzyme';
import Percent from './Percent';

import Enzyme from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';

Enzyme.configure({ adapter: new Adapter() });

test('Percent with 0', () => {
  const checkbox = shallow(<Percent value={0}/>);
  expect(checkbox.text()).toEqual('0%');
});

test('Percent with 0.10', () => {
  const checkbox = shallow(<Percent value={0.1}/>);
  expect(checkbox.text()).toEqual('10.00%');
});

test('Percent with 0.10', () => {
  const checkbox = shallow(<Percent value={0.1}/>);
  expect(checkbox.text()).toEqual('10.00%');
});

test('Percent with 0.16234', () => {
  const checkbox = shallow(<Percent value={0.16234}/>);
  expect(checkbox.text()).toEqual('16.23%');
});

test('Percent with 1.16234', () => {
  const checkbox = shallow(<Percent value={1.16234}/>);
  expect(checkbox.text()).toEqual('100%');
});

test('Percent with 1.00', () => {
  const checkbox = shallow(<Percent value={1.00}/>);
  expect(checkbox.text()).toEqual('100%');
});
