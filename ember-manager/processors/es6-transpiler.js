var path = process.argv[1],
    file = process.argv[2];

var Compiler = require("es6-module-transpiler").Compiler;

// TODO - path here should app's name instead of "app"
var module = path.match(/app\/(.+)\.[^\.]+$/)[1],
    output = (new Compiler(file, 'app/' + module)).toAMD();

process.stdout.write(output + '//# sourceURL=' + path + '\n');
