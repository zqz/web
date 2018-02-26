import FileSpeed from './FileSpeed';

it('works', () => {
  var now = (new Date()).getTime();

  const fs = new FileSpeed(100);

  fs.add(10, now);
  fs.add(10, now + 1000);
  fs.add(10, now + 2000);

  expect(fs.speed()).toEqual('15 Bytes/s');
});

it('works 2', () => {
  var now = (new Date()).getTime();

  const fs = new FileSpeed(100);

  fs.add(1024, now);
  fs.add(1024, now + 1000);
  fs.add(1024, now + 2000);

  expect(fs.speed()).toEqual('1.50 KB/s');
});
