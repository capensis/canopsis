import { setIn, unsetIn, addIn, removeIn } from '@/helpers/immutable';

import formComputedPropertiesMixin, { eventKeyComputed, formKeyComputed } from './internal/computed-properties';

/**
 * @mixin Form mixin
 */
export default {
  mixins: [formComputedPropertiesMixin],
  methods: {
    /**
     * Deep update field in object and in array
     *
     * @param {string|Array} path - Path to field or to array item
     * @param {*} value - New field or item value
     */
    updateField(path, value) {
      this.$emit(this[eventKeyComputed], setIn(this[this[formKeyComputed]], path, value));
    },

    /**
     * Deep remove field from object or item from array
     *
     * @param {string|Array} path - Path to field or to array item
     */
    removeField(path) {
      this.$emit(this[eventKeyComputed], unsetIn(this[this[formKeyComputed]], path));
    },

    /**
     * Deep add new item into array
     *
     * @param {string|Array} path - Path to field or to array item
     * @param {*} value - Value of new array item
     */
    addItemIntoArray(path, value) {
      this.$emit(this[eventKeyComputed], addIn(this[this[formKeyComputed]], path, value));
    },

    /**
     * Deep remove item from array
     *
     * @param {string|Array} path - Path to field or to array item
     * @param {number} index - Index of item
     */
    removeItemFromArray(path, index) {
      this.$emit(this[eventKeyComputed], removeIn(this[this[formKeyComputed]], path, index));
    },
  },
};
