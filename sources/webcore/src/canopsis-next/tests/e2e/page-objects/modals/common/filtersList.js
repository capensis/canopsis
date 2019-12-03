// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');
const el = require('../../../helpers/el');

const modalSelector = sel('filtersListModal');

const commands = {
  clickOutside() {
    return this.customClickOutside('@filtersList');
  },

  clickEditFilterByName(name) {
    return this.customClick(this.el('@editFilter', name));
  },

  verifyFilterVisibleByName(name) {
    return this.assert.visible(this.el('@filterItem', name));
  },

  verifyFilterDeletedByName(name) {
    return this.waitForElementNotPresent(this.el('@filterItem', name));
  },

  clickDeleteFilterByName(name) {
    return this.customClick(this.el('@deleteFilter', name));
  },

  clickAddFilter() {
    return this.customClick('@addFilter');
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      submitButton: sel('submitButton'),
    }),
    editFilter: `${sel('filtersListModal')} ${sel('editFilter-%s')}`,
    filterItem: `${sel('filtersListModal')} ${sel('filterItem-%s')}`,
    deleteFilter: `${sel('filtersListModal')} ${sel('deleteFilter-%s')}`,
    selectFilters: `${sel('filtersListModal')} ${sel('selectFilters')} .v-input__slot`,
    addFilter: `${sel('filtersListModal')} ${sel('addFilter')}`,

    filtersList: sel('filtersListModal'),
  },
  commands: [commands],
});
