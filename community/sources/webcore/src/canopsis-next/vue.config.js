// eslint-disable-next-line import/no-extraneous-dependencies
const TerserPlugin = require('terser-webpack-plugin');
const terserOptions = require('@vue/cli-service/lib/config/terserOptions');
const webpackOptions = require('@vue/cli-service/lib/options');

const updateFieldDirective = require('./tools/update-field-directive');

const isProduction = process.env.NODE_ENV === 'production';

module.exports = {
  publicPath: '/',
  lintOnSave: false,
  chainWebpack: (config) => {
    config.resolve.alias.store.set('vue$', 'vue/dist/vue.common.js');
    config.resolve.alias.store.set('handlebars', 'handlebars/dist/handlebars.js');

    config.module.rule('vue').use('vue-loader').loader('vue-loader')
      .tap((options) => {
        // eslint-disable-next-line no-param-reassign
        options.compilerOptions = {
          ...options.compilerOptions,

          directives: {
            field: updateFieldDirective,
          },
        };

        return options;
      })
      .end();

    config.optimization
      .minimizer('terser')
      .use(TerserPlugin, [{
        ...terserOptions(webpackOptions),
        exclude: /jodit/,
      }]);

    return config;
  },
  devServer: {
    host: 'localhost',
    https: true,
    disableHostCheck: true,
    watchOptions: {
      aggregateTimeout: 300,
      poll: true,
    },
  },
  pluginOptions: {
    webpackBundleAnalyzer: {
      analyzerMode: process.env.BUNDLE_ANALYZER_MODE, // 'disabled' / 'server' / 'static'
      openAnalyzer: false,
    },
    testAttrs: {
      enabled: isProduction,
      attrs: ['test'], // default: removes `data-test="..."`
    },
  },
};
