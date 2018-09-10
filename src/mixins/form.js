import uid from 'uid';

const eventKeyComputed = uid();

/**
 * @mixin Form mixin
 */
export default {
  data() {
    return {
      formKey: 'value',
    };
  },
  computed: {
    [eventKeyComputed]() {
      if (this.formKey === 'value') {
        return 'input';
      }

      return `update:${this.formKey}`;
    },
  },
  methods: {
    updateField(fieldName, value) {
      this.$emit(this[eventKeyComputed], { ...this[this.formKey], [fieldName]: value });
    },
    updateFieldInArrayItem(index, fieldName, value) {
      const items = [...this.columns];

      items[index] = {
        ...items[index],
        [fieldName]: this.prefixFormatter(value),
      };

      this.$emit(this[eventKeyComputed], items);
    },
  },
};
