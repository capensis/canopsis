import baseFormMixin, { modelPropKeyComputed, modelEventKeyComputed } from './base';

/**
 * @mixin Form mixin
 */
export default {
  mixins: [baseFormMixin],
  methods: {
    /**
     * Emit event to parent with new array with new item
     *
     * @param {*} value
     */
    addItemIntoArray(value) {
      this.$emit(this[modelEventKeyComputed], [...this[this[modelPropKeyComputed]], value]);
    },

    /**
     * Emit event to parent with new array with updated array item
     *
     * @param {number} index
     * @param {*} value
     */
    updateItemInArray(index, value) {
      const items = [...this[this[modelPropKeyComputed]]];

      items[index] = value;

      this.$emit(this[modelEventKeyComputed], items);
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
        ...this[this[modelPropKeyComputed]][index],
        [fieldName]: value,
      });
    },

    /**
     * Emit event to parent with new array without array item
     *
     * @param {number} index
     */
    removeItemFromArray(index) {
      this.$emit(this[modelEventKeyComputed], this[this[modelPropKeyComputed]].filter((v, i) => i !== index));
    },
  },
};
