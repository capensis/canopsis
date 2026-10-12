import { setField } from '@/helpers/immutable';

import { formBaseMixin, modelPropKeyComputed } from './base';

/**
 * @mixin Form mixin
 * @deprecated Should be used useModelField
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
  },
};

export default formMixin;
