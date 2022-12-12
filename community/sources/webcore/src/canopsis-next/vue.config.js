const MonacoWebpackPlugin = require('monaco-editor-webpack-plugin');

const updateFieldDirective = require('./tools/update-field-directive');

const isProduction = process.env.NODE_ENV === 'production';

module.exports = {
  publicPath: '/',
  lintOnSave: false,
  chainWebpack: (config) => {
    config.resolve.alias.store.set('vue$', 'vue/dist/vue.common.js');
    config.resolve.alias.store.set('handlebars', 'handlebars/dist/handlebars.js');

    config.plugin('monaco-editor-webpack-plugin')
      .use(MonacoWebpackPlugin, [{ languages: [] }]);

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
      .tap(([terserOptions]) => [{
        ...terserOptions,
        exclude: /TextEditor/,
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
