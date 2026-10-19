import { setField } from '@/helpers/immutable';

import { formBaseMixin, modelPropKeyComputed } from './base';

/**
 * @mixin Form mixin
 * @deprecated Should be used useArrayModelField
 */
export const formArrayMixin = {
  mixins: [formBaseMixin],
  methods: {
    /**
     * Emit event to parent with new array with new item
     *
     * @param {*} value
     * @return {Array}
     */
    addItemIntoArray(value) {
      return this.updateModel([...this[this[modelPropKeyComputed]], value]);
    },

    /**
     * Emit event to parent with new array with updated array item
     *
     * @param {number} index
     * @param {*} value
     * @return {Array}
     */
    updateItemInArray(index, value) {
      const items = [...this[this[modelPropKeyComputed]]];

      items[index] = value;

      return this.updateModel(items);
    },

    /**
     * Emit event to parent with new array with updated array item with updated field
     *
     * @param {number} index
     * @param {string} fieldName
     * @param {*} value
     * @return {Array}
     */
    updateFieldInArrayItem(index, fieldName, value) {
      return this.updateItemInArray(index, setField(this[this[modelPropKeyComputed]][index], fieldName, value));
    },

    /**
     * Emit event to parent with new array without array item
     *
     * @param {number} index
     * @return {Array}
     */
    removeItemFromArray(index) {
      return this.updateModel(this[this[modelPropKeyComputed]].filter((v, i) => i !== index));
    },
  },
};

export default formArrayMixin;
