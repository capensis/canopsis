// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  el,

  clickDropdownButton(tags) {
    return this.customClick(this.el('@dropdownButton', tags));
  },

  verifyDropdownZone(tags) {
    return this.waitForElementVisible(this.el('@dropdownZone', tags));
  },

  clickEditGroupButton(tags) {
    return this.waitForElementVisible(this.el('@editGroupButton', tags))
      .customClick(this.el('@editGroupButton', tags));
  },

  clickEditViewButton(title) {
    return this.waitForElementVisible(this.el('@editViewButton', title))
      .customClick(this.el('@editViewButton', title));
  },

  clickCopyViewButton(title) {
    return this.waitForElementVisible(this.el('@copyViewButton', title))
      .customClick(this.el('@copyViewButton', title));
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
};

module.exports = {
  elements: {
    userTopBarDropdownButton: sel('userTopBarDropdownButton'),
    userProfileButton: sel('userProfileButton'),
    logoutButton: sel('logoutButton'),
    dropdownButton: sel('dropButton-groupName-%s'),
    dropdownZone: `.v-menu__content ${sel('dropZone-groupName-%s')}`,
    editGroupButton: `${sel('dropButton-groupName-%s')} .v-btn${sel('editGroupButton')}`,
    editViewButton: `.v-btn${sel('editViewButton-viewTitle-%s')}`,
    copyViewButton: `.v-btn${sel('copyViewButton-viewTitle-%s')}`,
  },
  commands: [commands],
};
