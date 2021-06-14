// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  clickColorField() {
    return this.customClick('@chromeColorPickerField');
  },

  clearColorField() {
    return this.customClearValue('@chromeColorPickerField');
  },

  setColorField(value) {
    return this.customSetValue('@chromeColorPickerField', value);
  },
};

const modalSelector = sel('colorPickerModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      submitButton: sel('colorPickerSubmitButton'),
      cancelButton: sel('colorPickerCancelButton'),
    }),

    chromeColorPickerField: `${sel('colorPickerChrome')} .vc-editable-input input.vc-input__input`,
  },
  commands: [commands],
});
