// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../../helpers/el');
const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('infoPopupSettingModal');

const commands = {
  clickAddPopup() {
    return this.customClick('@addPopup');
  },

  clickEditPopup(index) {
    return this.customClick(this.el('@editPopup', index));
  },

  clickDeletePopup(index) {
    return this.customClick(this.el('@deletePopup', index));
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('infoPopupCancelButton'),
      submitButton: sel('infoPopupSubmitButton'),
    }),

    addPopup: sel('infoPopupAddPopup'),
    editPopup: `${sel('infoPopupSetting')}:nth-of-type(%s) ${sel('infoPopupEditPopup')}`,
    deletePopup: `${sel('infoPopupSetting')}:nth-of-type(%s) ${sel('infoPopupDeletePopup')}`,
  },
  commands: [commands],
});
