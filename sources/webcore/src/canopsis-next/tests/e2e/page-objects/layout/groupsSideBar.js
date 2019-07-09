// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');
const { elementsWrapperCreator } = require('../../helpers/page-object-creators');


const commands = {
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

  el,
};

const groupsSideBar = sel('groupsSideBar');

module.exports = {
  elements: {
    groupsSideBarButton: sel('groupsSideBarButton'),
    groupsSideBar,
    groupsSideBarClosed: `${groupsSideBar}.v-navigation-drawer--close`,

    ...elementsWrapperCreator(groupsSideBar, {
      groupSideBarSelectorByName: '.v-expansion-panel__header span[title="%s"]',
      groupSideBarSelectorById: `.v-expansion-panel__header ${sel('group-%s')}`,
      viewSideBarSelectorByName: '.v-expansion-panel__body a[title="%s"]',
      viewSideBarSelectorById: '.v-expansion-panel__body a[href^="/view/%s"]',
    }),
  },
  commands: [commands],
};
