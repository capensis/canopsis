<template>
  <v-tooltip location="left">
    <template #activator="{ props }">
      <v-badge :color="badgeColor" :content="count" class="logged-users-count">
        <v-btn size="x-small" variant="text" icon v-bind="props">
          <c-icon color="white" size="small"> people </c-icon>
        </v-btn>
      </v-badge>
    </template>

    <span>{{ $t('layout.sideBar.loggedUsersCount') }}</span>
  </v-tooltip>
</template>

<script>
import { mapActions } from 'pinia';

import { SOCKET_ROOMS } from '@/config';
import { useAuthStore } from '@/stores/modules/auth';

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

    this.$socket.join(SOCKET_ROOMS.loggedUserCount).addListener(this.setCount);
  },

  beforeUnmount() {
    this.$socket.leave(SOCKET_ROOMS.loggedUserCount).removeListener(this.setCount);
  },

  methods: {
    ...mapActions(useAuthStore, ['fetchLoggedUsersCountWithoutStore']),

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
    height: 17px;
    min-height: 17px;
    width: 17px;
    min-width: 17px;
    font-size: 12px;
  }
}
</style>
