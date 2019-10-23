const isProduction = process.env.NODE_ENV === 'production';

function isStringStart(chr) {
  const SINGLE_QUOTE_CODE = 0x22;
  const DOUBLE_QUOTE_CODE = 0x27;

  return chr === SINGLE_QUOTE_CODE || chr === DOUBLE_QUOTE_CODE;
}

function parseField(value) {
  const START_BRACKET_CODE = 0x5B;
  const END_BRACKET_CODE = 0x5D;
  const POINT_CODE = 0x2e;
  const path = [];
  let hasString = false;
  let isPoint = true;
  let startPos = 0;
  let endPos;

  for (let index = 0; index < value.length; index += 1) {
    const chr = value.charCodeAt(index);

    if (chr === START_BRACKET_CODE) {
      if (startPos < index - 1) {
        path.push(JSON.stringify(value.slice(startPos, index)));
      }

      hasString = false;
      isPoint = false;
      startPos = index + 1;

      if (isStringStart(value.charCodeAt(startPos))) {
        hasString = true;
      }
    } else if (chr === END_BRACKET_CODE) {
      endPos = index - 1;

      if (hasString) {
        path.push(JSON.stringify(value.slice(startPos + 1, endPos)));
        hasString = false;
      } else {
        path.push(value.slice(startPos, endPos + 1));
      }

      startPos = index + 1;
    } else if (chr === POINT_CODE && !hasString) {
      if (isPoint) {
        path.push(JSON.stringify(value.slice(startPos, index)));

        startPos = index + 1;
      }

      startPos = index + 1;
      isPoint = true;
    }
  }

  if (isPoint) {
    path.push(JSON.stringify(value.slice(startPos)));
  }

  return path;
}

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
            field(el, dir) {
              const { value, modifiers = {} } = dir;
              const { number, trim } = modifiers;
              const path = parseField(value.trim());

              path.shift();

              const baseValueExpression = '$$v';
              let valueExpression = baseValueExpression;

              if (trim) {
                valueExpression =
                  `(typeof ${baseValueExpression} === 'string'` +
                  `? ${baseValueExpression}.trim()` +
                  `: ${baseValueExpression})`;
              }

              if (number) {
                valueExpression = `_n(${valueExpression})`;
              }

              const assignment = `$updateField([${path}], ${valueExpression})`;

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
