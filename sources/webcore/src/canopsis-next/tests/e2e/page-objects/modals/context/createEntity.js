// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');
const el = require('../../../helpers/el');

const modalSelector = sel('createEntityModal');

const commands = {
  clickName() {
    return this.customClick('@entityName');
  },

  clearName() {
    return this.customClearValue('@entityName');
  },

  setName(value) {
    return this.customSetValue('@entityName', value);
  },

  setEnabled(checked = false) {
    return this.getAttribute('@entityEnabledInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@entityEnabled');
      }
    });
  },

  setType(index) {
    return this.customClick('@entitySelectType')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickImpact() {
    return this.customClick('@entityImpact');
  },

  clickDepends() {
    return this.customClick('@entityDepends');
  },

  clickTab(index) {
    return this.customClick(this.el('@createEntityTab', index));
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('createEntityCancelButton'),
      submitButton: sel('createEntitySubmitButton'),
    }),

    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',

    createEntityTab: `${sel('createEntityTab')}:nth-of-type(%s)`,
    entityName: `${sel('createEntityForm')} ${sel('entityFormName')}`,
    entityDescription: `${sel('createEntityForm')} ${sel('entityFormDescription')}`,
    entityEnabledInput: `${sel('createEntityForm')} input${sel('entityFormEnabled')}`,
    entityEnabled: `${sel('createEntityForm')} div${sel('entityFormEnabled')}`,
    entitySelectType: `${sel('createEntityForm')} div${sel('entityFormFieldLayout')} .v-select`,
    entityImpact: `${sel('createEntityForm')} ${sel('entityFormImpact')} ${sel('entitiesSelect')} .v-expansion-panel__header`,
    entityDepends: `${sel('createEntityForm')} ${sel('entityFormDepends')} ${sel('entitiesSelect')} .v-expansion-panel__header`,
  },
  commands: [commands],
});
