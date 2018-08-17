/**
 * @mixin Form mixin
 */
export default {
  data() {
    return {
      formKey: 'form',
    };
  },
  methods: {
    updateField(fieldName, value) {
      this.$emit(`update:${[this.formKey]}`, { ...this[this.formKey], [fieldName]: value });
    },
  },
};
