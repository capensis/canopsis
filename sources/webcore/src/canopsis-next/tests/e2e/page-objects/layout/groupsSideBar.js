// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');


const commands = {
  el,

  verifyPageElementsBefore() {
    return this.waitForElementVisible('@groupsSideBarButton');
  },

  clickPanelHeader(tags) {
    return this.customClick(this.el('@groupsSideBarHeader', tags));
  },

  verifyPanelBody(tags) {
    return this.waitForElementVisible(this.el('@groupsSideBarBody', tags));
  },

  clickGroupsSideBarButton() {
    this.api.element('css selector', this.el('@groupsSideBarClosed'), ({ status }) => {
      if (status === 0) {
        this.customClick('@groupsSideBarButton');
      }
    });

    return this;
  },

  clickEditGroupButton(tags) {
    return this.waitForElementVisible(this.el('@editGroupButton', tags))
      .customClick(this.el('@editGroupButton', tags));
  },

  clickEditViewButton(title) {
    return this.waitForElementVisible(this.el('@editViewButton', title))
      .customClick(this.el('@editViewButton', title));
  },

  clickLinkView(title) {
    return this.waitForElementVisible(this.el('@linkView', title))
      .customClick(this.el('@linkView', title));
  },

  clickCopyViewButton(title) {
    return this.waitForElementVisible(this.el('@copyViewButton', title))
      .customClick(this.el('@copyViewButton', title));
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
    groupsSideBarHeader: `${sel('panel-groupName-%s')} .v-expansion-panel__header`,
    groupsSideBarBody: `${sel('panel-groupName-%s')} .v-expansion-panel__body`,
    groupsSideBarButton: sel('groupsSideBarButton'),
    groupsSideBarClosed: `${groupsSideBar}.v-navigation-drawer--close`,
    groupSideBarSelectorById: `.v-expansion-panel__header ${sel('groupsSideBar-group-%s')}`,
    editGroupButton: `.v-expansion-panel__header .v-btn${sel('editGroupButton-groupName-%s')}`,
    editViewButton: `.v-expansion-panel__body .v-btn${sel('editViewButton-viewTitle-%s')}`,
    linkView: `.v-expansion-panel__body a.panel-item-content-link${sel('linkView-viewTitle-%s')}`,
    copyViewButton: `.v-expansion-panel__body .v-btn${sel('copyViewButton-viewTitle-%s')}`,
    viewSideBarSelectorById: '.v-expansion-panel__body a[href^="/view/%s"]',
  },
  commands: [commands],
};
