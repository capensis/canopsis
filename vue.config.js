module.exports = {
  lintOnSave: false,
  chainWebpack: (config) => {
    // config.output.publicPath('/en/static/cn/dist/');

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
