export const vuetifyTabsMixin = {
  methods: {
    callTabsOnResizeMethod() {
      this.$refs?.tabs?.onResize?.();
    },
    callTabsUpdateTabsMethod() {
      this.$refs?.tabs?.updateTabsView?.();
    },
  },
};
