import formComputedPropertiesMixin, { eventKeyComputed, formKeyComputed } from './internal/computed-properties';

/**
 * @mixin Form mixin
 */
export default {
  mixins: [formComputedPropertiesMixin],
  methods: {
    /**
     * Emit event to parent with new array with new item
     *
     * @param {*} value
     */
    addItemIntoArray(value) {
      this.$emit(this[eventKeyComputed], [...this[this[formKeyComputed]], value]);
    },

    /**
     * Emit event to parent with new array with updated array item
     *
     * @param {number} index
     * @param {*} value
     */
    updateItemInArray(index, value) {
      const items = [...this[this[formKeyComputed]]];

      items[index] = value;

      this.$emit(this[eventKeyComputed], items);
    },

    /**
     * Emit event to parent with new array with updated array item with updated field
     *
     * @param {number} index
     * @param {string} fieldName
     * @param {*} value
     */
    updateFieldInArrayItem(index, fieldName, value) {
      this.updateItemInArray(index, {
        ...this[this[formKeyComputed]][index],
        [fieldName]: value,
      });
    },

    /**
     * Emit event to parent with new array without array item
     *
     * @param {number} index
     */
    removeItemFromArray(index) {
      this.$emit(this[eventKeyComputed], this[this[formKeyComputed]].filter((v, i) => i !== index));
    },
  },
};
