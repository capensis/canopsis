const ELEMENTS_WAITING_DELAY = 5000;

module.exports.command = function(viewName, viewGroupname) {
  this
      .click('.v-toolbar__content .v-btn__content')
      .assert.containsText('.v-toolbar__content .v-btn__content', 'menu')
  return this;
};


