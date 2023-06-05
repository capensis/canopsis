<template lang="pug">
  v-tooltip(left)
    template(#activator="{ on }")
      v-badge.logged-users-count(:color="badgeColor", overlap)
        template(#badge="") {{ count }}
        v-btn(v-on="on", flat, icon, small)
          v-icon(color="white", small) people
    span {{ $t('layout.sideBar.loggedUsersCount') }}
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
  top: 0;

  .v-badge__badge {
    top: 2px;
    right: 2px;
    height: 17px;
    width: 17px;
    font-size: 12px;
  }
}
</style>
