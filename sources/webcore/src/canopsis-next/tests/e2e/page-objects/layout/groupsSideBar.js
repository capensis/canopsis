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

  clickEditGroupButton(id) {
    return this.customClick(this.el('@editGroupButton', id));
  },

  clickEditViewButton(id) {
    return this.customClick(this.el('@editViewButton', id));
  },

  clickCopyViewButton(id) {
    return this.customClick(this.el('@copyViewButton', id));
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
    editGroupButton: `.v-expansion-panel__header ${sel('groupsSideBar-editGroupButton-group-%s')}`,
    editViewButton: `.v-expansion-panel__header ${sel('groupsSideBar-editViewButton-group-%s')}`,
    copyViewButton: `.v-expansion-panel__header ${sel('groupsSideBar-copyViewButton-group-%s')}`,
    viewSideBarSelectorById: '.v-expansion-panel__body a[href^="/view/%s"]',
  },
  commands: [commands],
};
