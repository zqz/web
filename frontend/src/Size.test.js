import React from 'react';
import {shallow} from 'enzyme';
import Size from './Size';

import Enzyme from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';

Enzyme.configure({ adapter: new Adapter() });

test('Size with no bytes', () => {
  const checkbox = shallow(<Size bytes={0}/>);
  expect(checkbox.text()).toEqual('0');
});

test('Size with 1 bytes', () => {
  const checkbox = shallow(<Size bytes={1}/>);
  expect(checkbox.text()).toEqual('1 Bytes');
});

test('Size with 1KB', () => {
  const checkbox = shallow(<Size bytes={1024}/>);
  expect(checkbox.text()).toEqual('1.00 KB');
});

test('Size with 1MB', () => {
  const checkbox = shallow(<Size bytes={1024 * 1024}/>);
  expect(checkbox.text()).toEqual('1.00 MB');
});

test('Size with 1GB', () => {
  const checkbox = shallow(<Size bytes={1024 * 1024 * 1024}/>);
  expect(checkbox.text()).toEqual('1.00 GB');
});

test('Size with 1.5GB', () => {
  const checkbox = shallow(<Size bytes={1.5 * 1024 * 1024 * 1024}/>);
  expect(checkbox.text()).toEqual('1.50 GB');
});
