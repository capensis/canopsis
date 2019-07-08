// https://nightwatchjs.org/guide/#working-with-page-objects
const el = require('../helpers/el');
const { VUETIFY_ANIMATION_DELAY } = require('../../../src/config');

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

  clickOnEveryPopupsCloseIcons() {
    const { activePopupCloseIcon } = this.elements;

    return this.api.elements(activePopupCloseIcon.locateStrategy, activePopupCloseIcon.selector, ({ value = [] }) => {
      value.forEach(item => this.api.elementIdClick(item.ELEMENT).pause(VUETIFY_ANIMATION_DELAY));
    });
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

  el,
};

const groupSideBar = sel('groupsSideBar');

module.exports = {
  elements: {
    activePopupCloseIcon: `${sel('popupsWrapper')} .v-alert .v-alert__dismissible .v-icon`,
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
  },
  commands: [logoutPageCommands],
};
