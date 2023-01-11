export const validationAttachRequiredMixin = {
  methods: {
    attachRequiredRule(getter = () => this.value) {
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

    detachRequiredRule() {
      this.$validator.detach(this.name);
    },
  },
};
