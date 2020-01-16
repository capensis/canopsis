import { createNamespacedHelpers } from 'vuex';

import { DEFAULT_KEEP_ALIVE_INTERVAL } from '@/config';

import { getKeepalivePathByRoute } from '@/helpers/router';

const { mapActions } = createNamespacedHelpers('keepalive');

export default {
  data() {
    return {
      keepaliveTimer: undefined,
    };
  },
  methods: {
    ...mapActions(['keepalive']),

    async startKeepalive() {
      await this.keepalive({
        visible: !(document.visibilityState === 'hidden'),
        path: getKeepalivePathByRoute(this.$route),
      });

      if (!this.keepaliveTimer) {
        this.keepaliveTimer = setTimeout(this.startKeepalive, DEFAULT_KEEP_ALIVE_INTERVAL);
      }
    },

    stopKeepalive() {
      clearTimeout(this.keepaliveTimer);
      this.keepaliveTimer = undefined;
    },
  },
};
