const { process } = require('html-loader-jest');

module.exports = {
  process: content => ({
    code: process(content),
  }),
};
