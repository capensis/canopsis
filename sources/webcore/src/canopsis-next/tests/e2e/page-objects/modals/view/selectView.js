// https://nightwatchjs.org/guide/#working-with-page-objects

const { modalCreator } = require('../../../helpers/page-object-creators');
const el = require('../../../helpers/el');

const commands = {
  el,

  browseGroupById(id) {
    return this.customClick(this.el('@groupById', id));
  },

  browseViewById(id) {
    return this.customClick(this.el('@viewById', id));
  },
};

const modalSelector = sel('selectViewModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    groupById: sel('selectView-group-%s'),
    viewById: sel('selectView-view-%s'),
  },
  commands: [commands],
});
