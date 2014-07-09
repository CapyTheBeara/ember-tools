var fs = require('fs'),
    Compiler = require("es6-module-transpiler").Compiler,
    jsStringEscape = require('js-string-escape');

// https://github.com/joliss/broccoli-es6-concatenator
function wrapInEval (fileContents, fileName) {
  return 'eval("' +
    jsStringEscape(fileContents) +
    '//# sourceURL=' + jsStringEscape(fileName) +
    '");\n'
}

var path = process.argv[1],
    module = path.match(/app\/(.+)\.js/)[1];

try {
  var js = fs.readFileSync(path),
      output = (new Compiler(js, module)).toAMD();

  // TODO - path here should app's name instead of "app"
  process.stdout.write(wrapInEval(output, path));
} catch(e) {
  process.stderr.write('FILE_NOT_FOUND');
}



