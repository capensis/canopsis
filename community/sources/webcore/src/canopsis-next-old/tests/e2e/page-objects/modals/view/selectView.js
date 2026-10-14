// https://nightwatchjs.org/guide/#working-with-page-objects

const { modalCreator } = require('../../../helpers/page-object-creators');
const el = require('../../../helpers/el');

const commands = {
  el,

  browseGroupById(id) {
    return this.customClick(this.el('@groupById', id));
  },

  browseGroupByViewId(viewId) {
    return this.customClickXpath(this.el('@groupByViewId', viewId));
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

    groupByViewId: {
      locateStrategy: 'xpath',
      selector: '//div[contains(@class, "v-list__group")][div[div[a[contains(@data-test, "%s")]]]]//div[contains(@data-test, "selectView-group-")]',
    },
  },
  commands: [commands],
});
