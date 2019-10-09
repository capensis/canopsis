const isProduction = process.env.NODE_ENV === 'production';

module.exports = {
  baseUrl: isProduction ? '/en/static/canopsis-next/dist/' : '/',
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
            'field-model': function fieldModel(el, dir) {
              const { value } = dir;
              const path = value.split('.');

              path.shift();

              const baseValueExpression = '$$v';
              const assignment = `$updateField([${path.map(v => JSON.stringify(v)).join(', ')}], ${baseValueExpression})`;

              // eslint-disable-next-line no-param-reassign
              el.model = {
                value: `(${value})`,
                expression: JSON.stringify(value),
                callback: `function (${baseValueExpression}) {${assignment}}`,
              };
            },
          },
        };

        return options;
      });

    return config;
  },
  devServer: {
    proxy: {
      '/api': {
        target: process.env.VUE_APP_API_HOST,
        changeOrigin: true,
        pathRewrite: { '^/api': '' },
        secure: false,
        cookieDomainRewrite: '',
      },
      '/auth/external': {
        target: process.env.VUE_APP_API_HOST,
        changeOrigin: true,
        secure: false,
        cookieDomainRewrite: '',
      },
    },
    disableHostCheck: true,
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
