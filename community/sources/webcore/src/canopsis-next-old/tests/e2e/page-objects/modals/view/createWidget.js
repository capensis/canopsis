// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');
const el = require('../../../helpers/el');

const commands = {
  el,

  clickWidget(widget) {
    return this.customClick(this.el('@widget', widget));
  },
};

const modalSelector = sel('createWidgetModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      widget: sel('widget-%s'),
    }),
  },
  commands: [commands],
});
