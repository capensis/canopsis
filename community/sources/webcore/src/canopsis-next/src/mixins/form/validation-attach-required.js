/**
 * @deprecated Should be used useValidationAttachRequired
 */
export const validationAttachRequiredMixin = {
  methods: {
    attachRequiredRule(getter = () => this[this.$options?.model?.prop ?? 'value']) {
      const oldField = this.$validator.fields.find({ name: this.name });

      if (!oldField) {
        this.$validator.attach({
          name: this.name,
          rules: 'required:true',
          getter,
          vm: this,
        });
      }
    },

    validateRequiredRule() {
      this.$validator.validate(this.name);
    },

    detachRequiredRule() {
      this.$validator.detach(this.name);
    },
  },
};
