// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');
const el = require('../../../helpers/el');

const modalSelector = sel('pbehaviorListModal');

const commands = {
  clickAction(id, actionType) {
    return this.customClick(this.el('@alarmRowAction', id, actionType));
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      submitButton: sel('pbehaviorListConfirmButton'),
    }),

    alarmRowAction: sel('pbehaviorRow-%s-action-%s'),
  },
  commands: [commands],
});
