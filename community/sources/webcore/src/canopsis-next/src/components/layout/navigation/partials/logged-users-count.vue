<template>
  <v-tooltip left>
    <template #activator="{ on }">
      <v-badge
        class="logged-users-count"
        :color="badgeColor"
        overlap
      >
        <template #badge="">
          {{ count }}
        </template>
        <v-btn
          v-on="on"
          text
          icon
          small
        >
          <v-icon
            color="white"
            small
          >
            people
          </v-icon>
        </v-btn>
      </v-badge>
    </template>
    <span>{{ $t('layout.sideBar.loggedUsersCount') }}</span>
  </v-tooltip>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { SOCKET_ROOMS } from '@/config';

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
    this.fetchData();

    this.$socket
      .join(SOCKET_ROOMS.loggedUserCount)
      .addListener(this.setCount);
  },
  beforeDestroy() {
    this.$socket
      .leave(SOCKET_ROOMS.loggedUserCount)
      .removeListener(this.setCount);
  },
  methods: {
    ...mapActions(['fetchLoggedUsersCountWithoutStore']),

    setCount(count) {
      this.count = count;
    },

    async fetchData() {
      const { count = 0 } = await this.fetchLoggedUsersCountWithoutStore();

      this.setCount(count);
    },
  },
};
</script>

<style lang="scss">
.logged-users-count {
  position: absolute;
  margin: 6px;
  top: 0;

  .v-badge__badge {
    display: flex;
    align-items: center;
    justify-content: center;
    top: 2px;
    right: 2px;
    height: 17px;
    min-height: 17px;
    width: 17px;
    min-width: 17px;
    font-size: 12px;
  }
}
</style>
