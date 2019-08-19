const babelOptions = {
  presets: ['@vue/app'],
  plugins: ['require-context-hook', 'lodash'],
};

module.exports = require('babel-jest').createTransformer(babelOptions);
