// https://nightwatchjs.org/guide/#working-with-page-objects
const util = require('util');

const logoutPageCommands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@groupsSideBarButton');
  },

  clickGroupsSideBarButton() {
    this.api.element('css selector', this.el('@groupsSideBarClosed'), ({ status }) => {
      if (status === 0) {
        this.customClick('@groupsSideBarButton');
      }
    });

    return this;
  },

  browseGroupByName(name) {
    return this.customClick(this.el('@groupSideBarSelectorByName', name));
  },

  browseGroupById(id) {
    return this.customClick(this.el('@groupSideBarSelectorById', id));
  },

  browseViewByName(name) {
    return this.customClick(this.el('@viewSideBarSelectorByName', name));
  },

  browseViewById(id) {
    return this.customClick(this.el('@viewSideBarSelectorById', id));
  },

  el(elementName, data) {
    const element = this.elements[elementName.slice(1)];

    if (data) {
      return util.format(element.selector, data);
    }

    return element.selector;
  },
};

const groupSideBar = sel('groupsSideBar');

module.exports = {
  elements: {
    userTopBarDropdownButton: sel('userTopBarDropdownButton'),
    groupsSideBarButton: sel('groupsSideBarButton'),
    groupsSideBar: groupSideBar,
    groupsSideBarClosed: `${groupSideBar}.v-navigation-drawer--close`,
    groupSideBarSelectorByName: `${groupSideBar} .v-expansion-panel__header span[title="%s"]`,
    groupSideBarSelectorById: `${groupSideBar} .v-expansion-panel__header [data-test="group-%s"]`,
    viewSideBarSelectorByName: `${groupSideBar} .v-expansion-panel__body a[title="%s"]`,
    viewSideBarSelectorById: `${groupSideBar} .v-expansion-panel__body a[href^="/view/%s"]`,
    logoutButton: sel('logoutButton'),
    loginForm: sel('loginForm'),
    activePopup: `${sel('popupsWrapper')} .v-alert`,
  },
  commands: [logoutPageCommands],
};
