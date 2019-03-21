import { omit } from 'lodash';

import formComputedPropertiesMixin, { modelPropKeyComputed, modelEventKeyComputed } from './internal/computed-properties';

/**
 * @mixin Form mixin
 */
export default {
  mixins: [formComputedPropertiesMixin],
  methods: {
    /**
     * Emit event to parent with new object and with updated field
     *
     * @param {string} fieldName
     * @param {*} value
     */
    updateField(fieldName, value) {
      this.$emit(this[modelEventKeyComputed], { ...this[this[modelPropKeyComputed]], [fieldName]: value });
    },

    /**
     * Emit event to parent with new object
     * Rename a field in the object and update it
     *
     * @param {string} fieldName
     * @param {string} newFieldName
     * @param {*} value
     */
    updateAndMoveField(fieldName, newFieldName, value) {
      this.$emit(
        this[modelEventKeyComputed],
        { ...omit(this[this[modelPropKeyComputed]], fieldName), [newFieldName]: value },
      );
    },

    /**
     * Emit event to parent with new object without field
     *
     * @param {string} fieldName
     */
    removeField(fieldName) {
      this.$emit(this[modelEventKeyComputed], omit(this[this[modelPropKeyComputed]], [fieldName]));
    },
  },
};
