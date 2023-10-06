const MonacoWebpackPlugin = require('monaco-editor-webpack-plugin');

const updateFieldDirective = require('./tools/update-field-directive');

module.exports = {
  publicPath: '/',
  lintOnSave: false,
  chainWebpack: (config) => {
    config.resolve.alias.store.set('vue$', 'vue/dist/vue.common.js');
    config.resolve.alias.store.set('handlebars', 'handlebars/dist/handlebars.js');
    config.resolve.set(
      'fallback',
      {
        path: require.resolve('path-browserify'),
        process: require.resolve('process/browser'),
        url: require.resolve('url'),
      },
    );

    config.plugin('monaco-editor-webpack-plugin')
      .use(MonacoWebpackPlugin, [{ languages: [] }]);

    config.module.rule('html')
      .test(/^((?!index).)*\.html$/i)
      .use('html-loader')
      .loader('html-loader')
      .end();

    config.module.rule('vue')
      .use('vue-loader')
      .loader('vue-loader')
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
    allowedHosts: 'all',
    server: 'https',
    static: {
      watch: true,
    },
    client: {
      overlay: false,
    },
  },
};
