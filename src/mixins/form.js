/**
 * @mixin Form mixin
 */
export default {
  methods: {
    updateField(fieldName, value) {
      this.$emit('update:form', { ...this.form, [fieldName]: value });
    },
  },
};
