const { default: { createTransformer } } = require('babel-jest');

const babelOptions = {
  presets: ['@vue/app'],
  plugins: [
    '@babel/plugin-proposal-optional-chaining',
    '@babel/plugin-proposal-nullish-coalescing-operator',
    'require-context-hook',
    'lodash',
  ],
};

module.exports = createTransformer(babelOptions);
