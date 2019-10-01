// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitStatsNumber() {
    return this.customClick('@statsNumberSubmitButton');
  },

  clickDisplayModeEditButton() {
    return this.customClick('@displayModeEditButton');
  },

  clickSortOrder() {
    return this.customClick('@sortOrder');
  },

  selectSortOrder(index = 1) {
    return this.customClick('@sortOrderSelect')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  el,
};

module.exports = {
  elements: {
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s) .v-list__tile__content',

    statsNumberSubmitButton: sel('statsNumberSubmitButton'),

    displayModeEditButton: `${sel('statDisplayMode')} ${sel('editButton')}`,

    sortOrder: sel('sortOrder'),

    sortOrderSelect: `${sel('sortOrder')} .v-input__slot`,
  },
  commands: [commands],
};
