// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@page');
  },

  clearAppTitle() {
    return this.customClearValue('@appTitle');
  },

  setAppTitle(value) {
    return this.customSetValue('@appTitle', value);
  },

  selectLanguage(index = 1) {
    return this.customClick('@languageField')
      .waitForElementVisible(this.el('@languageOption', index))
      .customClick(this.el('@languageOption', index));
  },

  clearFooter() {
    return this.customClearRTE('@footerField');
  },

  clearDescription() {
    return this.customClearRTE('@descriptionField');
  },

  setFooter(value) {
    return this.customClick('@footerField')
      .sendKeys('@footerField', value);
  },

  setDescription(value) {
    return this.customClick('@descriptionField')
      .sendKeys('@descriptionField', value);
  },

  clickSubmitButton() {
    return this.customClick('@submitButton');
  },

  verifyAppTitle(text) {
    return this.assert.containsText('@appTitleElement', text);
  },

  verifyLoginDescription(text) {
    return this.assert.containsText('@loginDescriptionElement', text);
  },

  verifyLoginFooter(text) {
    return this.assert.containsText('@loginFooterElement', text);
  },

  el,
};

module.exports = {
  url() {
    return `${process.env.VUE_DEV_SERVER_URL}admin/parameters`;
  },
  elements: {
    page: sel('userInterfaceForm'),
    appTitle: sel('appTitle'),
    languageField: `${sel('languageLayout')} .v-input__slot`,
    languageOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    footerField: `${sel('footerLayout')} .jodit_wysiwyg`,
    descriptionFieldSource: `${sel('descriptionLayout')} .jodit_toolbar_btn`,
    descriptionField: `${sel('descriptionLayout')} .jodit_wysiwyg`,
    fileSelector: `input${sel('fileSelector')}`,
    appTitleElement: '#main-navigation .v-toolbar__title',
    loginDescriptionElement: '#main-content .mainContainer .description',
    loginFooterElement: `${sel('loginFormFooter')}`,
    submitButton: sel('submitButton'),
  },
  commands: [commands],
};
