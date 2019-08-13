// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  verifyPageBefore() {
    return this.waitForElementVisible('@adminRightsPage')
      .assert.visible('@addButton');
  },
  clickCheckbox(role, right, index) {
    return this.customClick(this.el('@checkbox', role, right, index))
      .defaultPause();
  },
  clickTab(tab) {
    return this.customClick(this.el('@tab', tab))
      .customClick(this.el('@tab', tab));
  },
  clickCreateUser() {
    return this.customClick('@createUser');
  },
  clickCreateRole() {
    return this.customClick('@createRole');
  },
  clickCreateRight() {
    return this.customClick('@createRight');
  },
  clickAddButton() {
    return this.customClick('@addButton');
  },
  clickSubmitRightButton() {
    return this.customClick('@submitRightButton');
  },
  el,
};

module.exports = {
  url() {
    return `${process.env.VUE_DEV_SERVER_URL}admin/rights`;
  },
  elements: {
    checkbox: `${sel('role-%s-right-%s')} .v-input:nth-child(%s) .v-input--selection-controls__ripple`,
    tab: `${sel('tab-%s')} a.v-tabs__item`,
    addButton: sel('addButton'),
    submitRightButton: sel('submitRightButton'),
    createUser: sel('createUser'),
    createRole: sel('createRole'),
    createRight: sel('createRight'),
    adminRightsPage: sel('adminRightsPage'),
  },
  commands: [commands],
};
