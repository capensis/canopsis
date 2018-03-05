var detect = require('feature-detect-es6')

if (!detect.newArrayFeatures()) {
  require('core-js/es6/array')
}

if (detect.all('class', 'arrowFunction')) {
  module.exports = require('./lib/command-line-commands')
} else {
  module.exports = require('./es5/command-line-commands')
}
