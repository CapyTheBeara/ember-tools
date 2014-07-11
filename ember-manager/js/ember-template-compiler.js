// var path = process.argv[1];
var compiler = require('ember-template-compiler')

var file = process.argv[2],
    output = compiler.precompile(file).toString(),
    template = "Ember.Handlebars.template(" + output + ");\n",
    es6 = "import Ember from 'ember';\nexport default " + template;

process.stdout.write(es6 + '\n');
