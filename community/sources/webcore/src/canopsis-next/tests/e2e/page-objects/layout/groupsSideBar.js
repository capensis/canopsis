// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  el,

  verifyPageElementsBefore() {
    return this.waitForElementVisible('@groupsSideBarButton');
  },

  clickPanelHeader(id) {
    return this.customClick(this.el('@groupsSideBarHeader', id));
  },

  verifyPanelBody(id) {
    return this.waitForElementVisible(this.el('@groupsSideBarBody', id));
  },

  groupsSideBarButtonElement(callback) {
    this.api.element('css selector', this.el('@groupsSideBarButton'), callback);

    return this;
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

  clickLinkView(id) {
    return this.waitForElementVisible(this.el('@linkView', id))
      .customClick(this.el('@linkView', id));
  },

  clickCopyViewButton(id) {
    return this.customClick(this.el('@copyViewButton', id));
  },

  browseGroupById(id) {
    this.api.element('css selector', this.el('@activeGroupSideBarSelectorById', id), ({ status }) => {
      if (status === -1) {
        this.customClick(this.el('@groupSideBarSelectorById', id));
      }
    });

    return this;
  },

  browseViewById(id) {
    return this.customClick(this.el('@viewSideBarSelectorById', id));
  },
};

const groupsSideBar = sel('groupsSideBar');

module.exports = {
  elements: {
    groupsSideBar,
    groupsSideBarHeader: `${sel('panel-group-%s')} .v-expansion-panel__header`,
    groupsSideBarBody: `${sel('panel-group-%s')} .v-expansion-panel__body`,
    groupsSideBarButton: sel('groupsSideBarButton'),
    groupsSideBarClosed: `${groupsSideBar}.v-navigation-drawer--close`,
    groupSideBarSelectorById: `.v-expansion-panel__header ${sel('groupsSideBar-group-%s')}`,
    linkView: `.v-expansion-panel__body a.panel-item-content-link${sel('linkView-view-%s')}`,
    activeGroupSideBarSelectorById: `.v-expansion-panel__container--active .v-expansion-panel__header ${sel('groupsSideBar-group-%s')}`,
    editGroupButton: `.v-expansion-panel__header .v-btn${sel('editGroupButton-group-%s')}`,
    editViewButton: `.v-expansion-panel__body .v-btn${sel('editViewButton-view-%s')}`,
    copyViewButton: `.v-expansion-panel__body .v-btn${sel('copyViewButton-view-%s')}`,
    viewSideBarSelectorById: '.v-expansion-panel__body a[href^="/view/%s"]',
  },
  commands: [commands],
};
