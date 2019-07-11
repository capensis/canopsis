// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');


const commands = {
  el,

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

  browseGroupById(id) {
    return this.customClick(this.el('@groupSideBarSelectorById', id));
  },

  browseViewById(id) {
    return this.customClick(this.el('@viewSideBarSelectorById', id));
  },
};

const groupsSideBar = sel('groupsSideBar');

module.exports = {
  elements: {
    groupsSideBar,

    groupsSideBarButton: sel('groupsSideBarButton'),
    groupsSideBarClosed: `${groupsSideBar}.v-navigation-drawer--close`,
    groupSideBarSelectorById: `.v-expansion-panel__header ${sel('groupsSideBar-group-%s')}`,
    viewSideBarSelectorById: '.v-expansion-panel__body a[href^="/view/%s"]',
  },
  commands: [commands],
};
