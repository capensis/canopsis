import isFunction from 'lodash/isFunction';

export default {
  callTabsOnResizeMethod() {
    if (this.$refs.tabs && isFunction(this.$refs.tabs.onResize)) {
      this.$refs.tabs.onResize();
    }
  },
};
