// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  verifyPageBefore() {
    return this.waitForElementVisible('@adminRightsPage')
      .assert.visible('@addButton');
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
  el,
};

module.exports = {
  url() {
    return `${process.env.VUE_DEV_SERVER_URL}admin/rights`;
  },
  elements: {
    tab: `${sel('tab-%s')} a.v-tabs__item`,
    addButton: sel('addButton'),
    createUser: sel('createUser'),
    createRole: sel('createRole'),
    createRight: sel('createRight'),
    adminRightsPage: sel('adminRightsPage'),
  },
  commands: [commands],
};
