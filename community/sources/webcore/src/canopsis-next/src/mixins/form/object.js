import { isEqual } from 'lodash';

import { setField, unsetField } from '@/helpers/immutable';

import { formBaseMixin, modelPropKeyComputed } from './base';

/**
 * @mixin Form mixin
 */
export const formMixin = {
  mixins: [formBaseMixin],
  methods: {
    /**
     * Emit event to parent with new object and with updated field
     *
     * @param {string|Array} fieldName
     * @param {*} value
     */
    updateField(fieldName, value) {
      this.updateModel(setField(this[this[modelPropKeyComputed]], fieldName, value));
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
        const result = unsetField(this[this[modelPropKeyComputed]], fieldName);

        this.updateModel(setField(result, newFieldName, value));
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
        const result = unsetField(this[this[modelPropKeyComputed]], fieldName);

        this.updateModel(setField(result, newFieldName, value => value));
      }
    },

    /**
     * Emit event to parent with new object without field
     *
     * @param {string|Array} fieldName
     */
    removeField(fieldName) {
      this.updateModel(unsetField(this[this[modelPropKeyComputed]], fieldName));
    },
  },
};

export default formMixin;
