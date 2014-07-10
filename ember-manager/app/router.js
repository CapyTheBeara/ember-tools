import Ember from 'ember';

var Router = Ember.Router.extend({
  location: 'hash'
});

Router.map(function() {
  this.route('foo');
});

export default Router;
