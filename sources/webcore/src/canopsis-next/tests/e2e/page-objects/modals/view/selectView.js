// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');
const el = require('../../../helpers/el');

const commands = {
  browseGroupById(id) {
    return this.customClick(this.el('@groupById', id));
  },

  browseViewById(id) {
    return this.customClick(this.el('@viewById', id));
  },

  el,
};

const modalSelector = sel('selectViewModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      groupById: sel('group-%s'),
      viewById: sel('view-%s'),
    }),
  },
  commands: [commands],
});
