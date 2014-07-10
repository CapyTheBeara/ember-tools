var fs = require('fs'),
    Compiler = require("es6-module-transpiler").Compiler;

function addSourceURL (fileContents, fileName) {
  return fileContents + '//# sourceURL=' + fileName;
}

// TODO - path here should app's name instead of "app"
var path = process.argv[1],
    module = path.match(/app\/(.+)\.js/)[1];

try {
  var js = fs.readFileSync(path),
      output = (new Compiler(js, 'app/' + module)).toAMD();

  process.stdout.write(addSourceURL(output, path));
} catch(e) {
  process.stderr.write('FILE_NOT_FOUND');
}



