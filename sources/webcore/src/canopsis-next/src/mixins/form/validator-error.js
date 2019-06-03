export default {
  computed: {
    hasAnyError() {
      if (this.$validator && this.errors.any()) {
        const fields = this.$validator.fields.items.filter(({ vm }) => this.isParent(vm));

        return fields.some(field => this.errors.has(field.name));
      }

      return false;
    },
  },
  methods: {
    isParent(vm) {
      if (vm) {
        if (vm === this || vm._original === this) {
          return true;
        }

        if (vm.$parent) {
          return this.isParent(vm.$parent);
        }
      }

      return false;
    },
  },
};
