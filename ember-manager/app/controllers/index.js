function foo() {
  return 'bar';
}

export default Em.Controller.extend({
  world: 'world',
  foo: foo.property()
});

