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
  setFooter(value) {
    return this.url('https://xdsoft.net/jodit/')
      .customClick('@footerField')
      .sendKeys('@footerField', value)
      .pause(50000);
  },
  setDescription(value) {
    return this.url('https://xdsoft.net/jodit/')
      .customClick('@descriptionField')
      .sendKeys('@descriptionField', value)
      .pause(50000);
  },
  clickSubmitButton() {
    return this.customClick('@submitButton');
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
    descriptionField: `${sel('descriptionLayout')} .jodit_wysiwyg`,
    fileSelector: `input${sel('fileSelector')}`,
    submitButton: sel('submitButton'),
  },
  commands: [commands],
};
