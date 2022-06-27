import { isFunction } from 'lodash';

export const vuetifyTabsMixin = {
  methods: {
    callTabsOnResizeMethod() {
      if (this.$refs.tabs && isFunction(this.$refs.tabs.onResize)) {
        this.$refs.tabs.onResize();
      }
    },
    callTabsUpdateTabsMethod() {
      if (this.$refs.tabs && isFunction(this.$refs.tabs.updateTabsView)) {
        this.$refs.tabs.updateTabsView();
      }
    },
  },
};
