<template lang="pug">
  v-tooltip.active-sessions-count(left)
    v-badge(slot="activator", :color="badgeColor", right, overlap)
      span(slot="badge") {{ count }}
      v-btn(flat, icon, small)
        v-icon(color="white", small) people
    span {{ $t('layout.sideBar.activeSessions') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { ACTIVE_SESSIONS_COUNT_FETCHING_INTERVAL } from '@/config';

const { mapActions } = createNamespacedHelpers('session');

export default {
  props: {
    badgeColor: {
      type: String,
      default: 'primary',
    },
  },
  data() {
    return {
      count: '',
      requestTimer: undefined,
    };
  },
  mounted() {
    this.startFetchActiveSessionsCount();
  },
  beforeDestroy() {
    this.stopFetchActiveSessionsCount();
  },
  methods: {
    ...mapActions({
      fetchSessionsListWithoutStore: 'fetchListWithoutStore',
    }),

    async startFetchActiveSessionsCount() {
      const { sessions } = await this.fetchSessionsListWithoutStore({ params: { active: true } });

      this.count = sessions.length;

      if (this.requestTimer) {
        this.stopFetchActiveSessionsCount();
      }

      this.requestTimer = setTimeout(this.startFetchActiveSessionsCount, ACTIVE_SESSIONS_COUNT_FETCHING_INTERVAL);
    },

    stopFetchActiveSessionsCount() {
      clearTimeout(this.requestTimer);

      this.requestTimer = undefined;
    },
  },
};
</script>

<style lang="scss" scoped>
  .active-sessions-count /deep/ .v-badge__badge {
    top: 2px;
    right: 2px;
    height: 17px;
    width: 17px;
    font-size: 12px;
  }
</style>
