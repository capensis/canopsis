// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
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
      cancelButton: sel('infoPopupCancelButton'),
      submitButton: sel('infoPopupSubmitButton'),

      addPopup: sel('infoPopupAddPopup'),
      editPopup: sel('infoPopupEditPopup'),
      deletePopup: sel('infoPopupDeletePopup')
    }),
  },
  commands: [commands],
});
