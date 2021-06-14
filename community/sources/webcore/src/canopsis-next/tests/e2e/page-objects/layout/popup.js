// https://nightwatchjs.org/guide/#working-with-page-objects
const { VUETIFY_ANIMATION_DELAY } = require('../../../../src/config');

const commands = {
  clickOnEveryPopupsCloseIcons() {
    const { activePopupCloseIcon } = this.elements;

    return this.api.elements(activePopupCloseIcon.locateStrategy, activePopupCloseIcon.selector, ({ value = [] }) => {
      value.forEach(item => this.api.elementIdClick(item.ELEMENT).pause(VUETIFY_ANIMATION_DELAY));
    });
  },
};

module.exports = {
  elements: {
    activePopupCloseIcon: `${sel('popupsWrapper')} .v-alert .v-alert__dismissible .v-icon`,
  },
  commands: [commands],
};
