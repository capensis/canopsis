import { createNamespacedHelpers } from 'vuex';

import { DEFAULT_KEEP_ALIVE_INTERVAL } from '@/config';

const { mapActions } = createNamespacedHelpers('keepalive');

export default {
  methods: {
    ...mapActions(['keepAlive', 'sessionHide']),

    startKeepAlive() {
      if (this.keepAliveInterval === undefined) {
        this.keepAliveInterval = setInterval(() => {
          this.keepAlive({
            visible: this.getvisible(),
            path: this.getPath(),
          });
        }, DEFAULT_KEEP_ALIVE_INTERVAL);
      }
    },
    stopKeepAlive() {
      clearInterval(this.keepAliveInterval);
      this.keepAliveInterval = undefined;
    },
    startSessionHide() {
      this.sessionHide({
        path: this.getPath(),
      });
    },
    getvisible() {
      if (document.visibilityState === 'hidden') {
        return false;
      }
      return true;
    },
    getPath() {
      const { tabId } = this.$route.query;
      if (tabId) {
        return [this.$route.path, tabId];
      }
      return [this.$route.path];
    },
  },
};
