const { createTransformer } = require('babel-jest');

const babelOptions = {
  presets: ['@vue/app'],
  plugins: ['require-context-hook', 'lodash'],
};

module.exports = createTransformer(babelOptions);
