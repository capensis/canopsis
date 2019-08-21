// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  clickSubmitButton() {
    return this.customClick('@submitButton');
  },
  clickAddPopup() {
    return this.customClick('@addPopup');
  },
  clickEditPopup() {
    return this.customClick('@editPopup');
  },
  clickDeletePopup() {
    return this.customClick('@deletePopup');
  },
};

const modalSelector = sel('infoPopupSetting');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      addPopup: sel('addPopup'),
      editPopup: sel('editPopup'),
      deletePopup: sel('deletePopup'),
      submitButton: sel('submitButton'),
    }),
  },
  commands: [commands],
});
