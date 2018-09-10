import uid from '@/helpers/uid';

const eventKeyComputed = uid('eventKey');
const formKeyComputed = uid('formKey');

/**
 * @mixin Form mixin
 */
const formMixin = {
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

export default formMixin;
