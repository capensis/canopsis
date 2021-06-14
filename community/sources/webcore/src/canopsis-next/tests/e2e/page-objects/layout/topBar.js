// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  el,

  clickDropdownButton(id) {
    return this.customClick(this.el('@dropdownButton', id));
  },

  verifyDropdownZone(id) {
    return this.waitForElementVisible(this.el('@dropdownZone', id));
  },

  clickEditGroupButton(id) {
    return this.customClick(this.el('@editGroupButton', id));
  },

  clickEditViewButton(id) {
    return this.customClick(this.el('@editViewButton', id));
  },

  clickCopyViewButton(id) {
    return this.customClick(this.el('@copyViewButton', id));
  },

  clickUserDropdown() {
    return this.customClick('@userTopBarDropdownButton');
  },

  clickUserProfileButton() {
    return this.customClick('@userProfileButton');
  },

  clickLogoutButton() {
    return this.customClick('@logoutButton');
  },

  getLogoutButtonText(callback) {
    return this.getText('@logoutButtonText', callback);
  },
};

const logoutButtonSelector = sel('logoutButton');

module.exports = {
  elements: {
    userTopBarDropdownButton: sel('userTopBarDropdownButton'),
    userProfileButton: sel('userProfileButton'),
    logoutButton: logoutButtonSelector,
    logoutButtonText: `${logoutButtonSelector} .ml-2`,
    dropdownButton: sel('dropDownButton-group-%s'),
    dropdownZone: `.v-menu__content ${sel('dropDownZone-group-%s')}`,
    editGroupButton: `${sel('dropDownButton-group-%s')} .v-btn${sel('editGroupButton')}`,
    editViewButton: `.v-btn${sel('editViewButton-view-%s')}`,
    copyViewButton: `.v-btn${sel('copyViewButton-view-%s')}`,
  },
  commands: [commands],
};
