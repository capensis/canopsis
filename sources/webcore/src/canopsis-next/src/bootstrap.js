import WebFont from 'webfontloader';

WebFont.load({
  custom: {
    families: ['Roboto:300,400,500,700', 'Material+Icons'],
    urls: process.env.NODE_ENV === 'production' ? ['./styles/fonts.css'] : ['/styles/fonts.css'],
  },
});
