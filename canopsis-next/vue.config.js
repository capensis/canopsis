module.exports = {
  lintOnSave: false,
  chainWebpack: (config) => {
    if (process.env.NODE_ENV === 'production') {
      config.output.publicPath('/en/static/canopsis-next/dist/');
    }
    config.resolve.alias.store.set('handlebars', 'handlebars/dist/handlebars.js');

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
    },
  },
};
