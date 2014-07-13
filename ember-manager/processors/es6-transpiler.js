var path = require('path'),
    Compiler = require("es6-module-transpiler").Compiler;

var filePath = process.argv[1],
    file = process.argv[2];

var projectName = path.basename(process.cwd()),
    module = filePath.match(/app\/(.+)\.[^\.]+$/)[1],
    newPath = path.join(projectName, module),
    output = (new Compiler(file, newPath)).toAMD();

process.stdout.write(output + '//# sourceURL=' + newPath + '.js' + '\n');
