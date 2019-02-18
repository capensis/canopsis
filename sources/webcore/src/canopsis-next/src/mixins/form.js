import { omit } from 'lodash';

import uid from '@/helpers/uid';
import { setIn, unsetIn, addIn, removeIn } from '@/helpers/immutable';

const eventKeyComputed = uid('_eventKey');
const formKeyComputed = uid('_formKey');

/**
 * @mixin Form mixin
 */
export default {
  computed: {
    [formKeyComputed]() {
      if (this.$options.model && this.$options.model.prop) {
        return this.$options.model.prop;
      }

      return 'value';
    },
    [eventKeyComputed]() {
      if (this.$options.model && this.$options.model.event) {
        return this.$options.model.event;
      }

      return 'input';
    },
  },
  methods: {
    /**
     * Emit event to parent with new object and with updated field
     *
     * @param {string} fieldName
     * @param {*} value
     */
    updateField(fieldName, value) {
      this.$emit(this[eventKeyComputed], { ...this[this[formKeyComputed]], [fieldName]: value });
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
        this[eventKeyComputed],
        { ...omit(this[this[formKeyComputed]], fieldName), [newFieldName]: value },
      );
    },

    /**
     * Emit event to parent with new object without field
     *
     * @param {string} fieldName
     */
    removeField(fieldName) {
      this.$emit(this[eventKeyComputed], omit(this[this[formKeyComputed]], [fieldName]));
    },

    /**
     * Emit event to parent with new array with new item
     *
     * @param {*} value
     */
    addItemIntoArray(value) {
      this.$emit(this[eventKeyComputed], [...this[this[formKeyComputed]], value]);
    },

    /**
     * Emit event to parent with new array with updated array item with updated field
     *
     * @param {number} index
     * @param {string} fieldName
     * @param {*} value
     */
    updateFieldInArrayItem(index, fieldName, value) {
      const items = [...this[this[formKeyComputed]]];

      items[index] = {
        ...items[index],
        [fieldName]: value,
      };

      this.$emit(this[eventKeyComputed], items);
    },

    /**
     * Emit event to parent with new array without array item
     *
     * @param {number} index
     */
    removeItemFromArray(index) {
      this.$emit(this[eventKeyComputed], this[this[formKeyComputed]].filter((v, i) => i !== index));
    },

    /**
     * Deep update field in object and in array
     *
     * @param {string|Array} path - Path to field or to array item
     * @param {*} value - New field or item value
     */
    deepUpdateField(path, value) {
      this.$emit(this[eventKeyComputed], setIn(this[this[formKeyComputed]], path, value));
    },

    /**
     * Deep remove field from object or item from array
     *
     * @param {string|Array} path - Path to field or to array item
     */
    deepRemoveField(path) {
      this.$emit(this[eventKeyComputed], unsetIn(this[this[formKeyComputed]], path));
    },

    /**
     * Deep add new item into array
     *
     * @param {string|Array} path - Path to field or to array item
     * @param {*} value - Value of new array item
     */
    deepAddItemIntoArray(path, value) {
      this.$emit(this[eventKeyComputed], addIn(this[this[formKeyComputed]], path, value));
    },

    /**
     * Deep remove item from array
     *
     * @param {string|Array} path - Path to field or to array item
     * @param {number} index - Index of item
     */
    deepRemoveItemFromArray(path, index) {
      this.$emit(this[eventKeyComputed], removeIn(this[this[formKeyComputed]], path, index));
    },
  },
};
