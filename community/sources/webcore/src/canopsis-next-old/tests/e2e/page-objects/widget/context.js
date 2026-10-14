// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitContext() {
    return this.customClick('@submitContext');
  },

  clickContextTypeOfEntities() {
    return this.customClick('@contextTypeOfEntities');
  },

  selectEntitiesTypeCheckbox(index, checked) {
    return this.getAttribute(
      this.el('@entitiesTypeCheckboxInput', index),
      'aria-checked',
      ({ value }) => {
        if (value !== String(checked)) {
          this.customClick(this.el('@entitiesTypeCheckbox', index));
        }
      },
    );
  },

  el,
};

module.exports = {
  elements: {
    submitContext: sel('submitContext'),

    contextTypeOfEntities: sel('contextTypeOfEntities'),

    entitiesTypeCheckbox: `div${sel('entitiesTypeCheckbox')}:nth-of-type(%s) .v-input--selection-controls__ripple`,
    entitiesTypeCheckboxInput: `div${sel('entitiesTypeCheckbox')}:nth-of-type(%s) input${sel('entitiesTypeCheckbox')}`,
  },
  commands: [commands],
};
