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
            visible: !(document.visibilityState === 'hidden'),
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
    getPath() {
      const { tabId } = this.$route.query;
      if (tabId) {
        return [this.$route.path, tabId];
      }
      return [this.$route.path];
    },
  },
};
