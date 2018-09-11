import omit from 'lodash/omit';

import uid from '@/helpers/uid';

const eventKeyComputed = uid('eventKey');
const formKeyComputed = uid('formKey');

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
    updateField(fieldName, value) {
      this.$emit(this[eventKeyComputed], { ...this[this[formKeyComputed]], [fieldName]: value });
    },
    deleteField(fieldName) {
      this.$emit(this[eventKeyComputed], omit(this[this[formKeyComputed]], [fieldName]));
    },
    updateFieldInArrayItem(index, fieldName, value) {
      const items = [...this[this[formKeyComputed]]];

      items[index] = {
        ...items[index],
        [fieldName]: value,
      };

      this.$emit(this[eventKeyComputed], items);
    },
  },
};
