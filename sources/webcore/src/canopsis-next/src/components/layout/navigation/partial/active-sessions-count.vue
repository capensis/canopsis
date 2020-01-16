<template lang="pug">
  v-tooltip.active-sessions-count(left)
    v-badge(slot="activator", right, overlap)
      span(slot="badge") {{ count }}
      v-btn(flat, icon, small)
        v-icon(color="white", small) people
    span Active sessions
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { ACTIVE_SESSION_FETCHING_INTERVAL } from '@/config';

const { mapActions } = createNamespacedHelpers('session');

export default {
  data() {
    return {
      count: '',
      requestTimer: undefined,
    };
  },
  mounted() {
    this.fetchActiveSessionsCount();
  },
  beforeDestroy() {
    if (this.requestTimer) {
      clearTimeout(this.requestTimer);
      this.requestTimer = undefined;
    }
  },
  methods: {
    ...mapActions({
      fetchSessionsListWithoutStore: 'fetchListWithoutStore',
    }),

    async fetchActiveSessionsCount() {
      const { data: sessions } = await this.fetchSessionsListWithoutStore({ active: true });

      this.count = sessions.length;

      if (!this.requestTimer) {
        this.requestTimer = setTimeout(this.fetchActiveSessionsCount, ACTIVE_SESSION_FETCHING_INTERVAL);
      }
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
