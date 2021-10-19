<template lang="pug">
  v-tooltip.logged-users-count(left)
    v-badge(slot="activator", :color="badgeColor", right, overlap)
      span(slot="badge") {{ count }}
      v-btn(flat, icon, small)
        v-icon(color="white", small) people
    span {{ $t('layout.sideBar.loggedUsersCount') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { ACTIVE_LOGGED_USERS_COUNT_FETCHING_INTERVAL } from '@/config';

const { mapActions } = createNamespacedHelpers('auth');

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
    };
  },
  mounted() {
    this.startFetchLoggedUsersCount();
  },
  beforeDestroy() {
    this.stopFetchLoggedUsersCount();
  },
  methods: {
    ...mapActions(['fetchLoggedUsersCountWithoutStore']),

    async startFetchLoggedUsersCount() {
      const { count } = await this.fetchLoggedUsersCountWithoutStore();

      this.count = count;

      if (this.requestTimer) {
        this.stopFetchLoggedUsersCount();
      }

      this.requestTimer = setTimeout(this.startFetchLoggedUsersCount, ACTIVE_LOGGED_USERS_COUNT_FETCHING_INTERVAL);
    },

    stopFetchLoggedUsersCount() {
      clearTimeout(this.requestTimer);

      this.requestTimer = undefined;
    },
  },
};
</script>

<style lang="scss" scoped>
  .logged-users-count /deep/ .v-badge__badge {
    top: 2px;
    right: 2px;
    height: 17px;
    width: 17px;
    font-size: 12px;
  }
</style>
