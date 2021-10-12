import { createNamespacedHelpers } from 'vuex';

import { getViewStatsPathByRoute } from '@/helpers/router';
import { DEFAULT_VIEW_STATS_INTERVAL } from '@/config';

const { mapActions } = createNamespacedHelpers('viewStats');

export const entitiesViewStatsMixin = {
  methods: {
    ...mapActions({
      updateViewStats: 'update',
    }),

    async startViewStats() {
      await this.updateViewStats({
        data: {
          visible: !(document.visibilityState === 'hidden'),
          path: getViewStatsPathByRoute(this.$route),
        },
      });

      this.stopViewStats();

      this.requestTimer = setTimeout(this.startViewStats, DEFAULT_VIEW_STATS_INTERVAL);
    },

    stopViewStats() {
      clearTimeout(this.requestTimer);
    },
  },
};
