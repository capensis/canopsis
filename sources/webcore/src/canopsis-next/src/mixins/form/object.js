import { isEqual } from 'lodash';

import { setIn, unsetIn } from '@/helpers/immutable';

import baseFormMixin, { modelPropKeyComputed } from './base';

/**
 * @mixin Form mixin
 */
export default {
  mixins: [baseFormMixin],
  methods: {
    /**
     * Emit event to parent with new object and with updated field
     *
     * @param {string|Array} fieldName
     * @param {*} value
     */
    updateField(fieldName, value) {
      this.updateModel(setIn(this[this[modelPropKeyComputed]], fieldName, value));
    },

    /**
     * Emit event to parent with new object
     * Rename a field in the object and update it
     *
     * @param {string|Array} fieldName
     * @param {string|Array} newFieldName
     * @param {*} value
     */
    updateAndMoveField(fieldName, newFieldName, value) {
      if (isEqual(fieldName, newFieldName)) {
        this.updateField(fieldName, value);
      } else {
        const result = unsetIn(this[this[modelPropKeyComputed]], fieldName);

        this.updateModel(setIn(result, newFieldName, value));
      }
    },
    /**
     * Emit event to parent with new object
     * Rename a field in the object
     *
     * @param {string|Array} fieldName
     * @param {string|Array} newFieldName
     */
    moveField(fieldName, newFieldName) {
      if (!isEqual(fieldName, newFieldName)) {
        const result = unsetIn(this[this[modelPropKeyComputed]], fieldName);

        this.updateModel(setIn(result, newFieldName, value => value));
      }
    },

    /**
     * Emit event to parent with new object without field
     *
     * @param {string|Array} fieldName
     */
    removeField(fieldName) {
      this.updateModel(unsetIn(this[this[modelPropKeyComputed]], fieldName));
    },
  },
};
