
module.exports = {
  '@tags': ['manageviews'],
  'Homepage and login' : function (browser) {
    browser
       .login('root', 'root')
  },
  'Click on menu' : function (browser) {
    browser
      .click('.v-toolbar__content .v-btn__content')
      .end();
  }
};

