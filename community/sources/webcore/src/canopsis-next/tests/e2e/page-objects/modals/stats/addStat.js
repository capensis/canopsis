// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../../helpers/el');
const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('addStatModal');

const commands = {
  selectStatType(index = 1) {
    return this.customClick('@statType')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickStatTitle() {
    return this.customClick('@statTitle');
  },

  clearStatTitle() {
    return this.customClearValue('@statTitle');
  },

  setStatTitle(value) {
    return this.customSetValue('@statTitle', value);
  },

  setStatTrend(checked) {
    return this.getAttribute('@statTrendInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@statTrend');
      }
    });
  },

  setStatRecursive(checked) {
    return this.getAttribute('@statRecursiveInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@statRecursive');
      }
    });
  },

  clickStatStates() {
    return this.customClick('@statStates');
  },

  setStatState(index, checked) {
    return this.getAttribute(
      this.el('@statStatesOptionInput', index),
      'aria-checked',
      ({ value }) => {
        if (value !== String(checked)) {
          this.customClick(this.el('@statStatesOption', index));
        }
      },
    );
  },

  setStatStates(states = []) {
    states.forEach(({ index, checked }) => {
      this.setStatState(index, checked);
    });
    return this;
  },

  clickStatAuthors() {
    return this.customClick('@statAuthors');
  },

  clearStatAuthors() {
    return this.customClearValue('@statAuthors');
  },

  setStatAuthor(value) {
    return this.customSetValue('@statAuthors', value)
      .customKeyup('@statAuthors', this.api.Keys.ENTER);
  },

  setStatAuthors(authors) {
    authors.forEach((author) => {
      this.setStatAuthor(author);
    });
    return this;
  },

  removeAuthor(value) {
    this.api
      .useXpath()
      .customClick(this.el('@statAuthorXPath', value));

    this.sendKeys('@statAuthors', this.api.Keys.BACK_SPACE);

    return this;
  },

  clickParameters() {
    return this.customClick('@statParameters');
  },

  clickStatSla() {
    return this.customClick('@statSla');
  },

  clearStatSla() {
    return this.customClearValue('@statSla');
  },

  setStatSla(value) {
    return this.customSetValue('@statSla', value);
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('addStatCancelButton'),
      submitButton: sel('addStatSubmitButton'),
    }),
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s) .v-list__tile__content',

    addStatModal: sel('addStatModal'),

    statType: `${sel('addStatModal')} div${sel('statTypeLayout')} .v-input__slot`,

    statTitle: `${sel('addStatModal')} ${sel('statTitle')}`,

    statSla: `${sel('addStatModal')} ${sel('statSla')}`,

    statTrend: `${sel('addStatModal')} div${sel('statTrend')} .v-input__slot`,
    statTrendInput: `${sel('addStatModal')} input${sel('statTrend')}`,

    statRecursive: `${sel('addStatModal')} div${sel('statRecursive')} .v-input__slot`,
    statRecursiveInput: `${sel('addStatModal')} input${sel('statRecursive')}`,

    statStates: `${sel('addStatModal')} ${sel('statStates')} .v-input__slot`,
    statStatesOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    statStatesOptionInput: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s) input',

    statAuthors: `${sel('addStatModal')} input${sel('statAuthors')}`,
    statAuthorXPath: './/*[@data-test=\'addStatModal\']//div[@class=\'v-select__selections\']//span[span[@class=\'v-chip__content\' and contains(text(), \'%s\')]]',

    statParameters: sel('statParameters'),
  },
  commands: [commands],
});
