<template lang="pug">
  div
    v-btn(@click="click") Open watcher modal
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import modalMixin from '@/mixins/modal/modal';
import watcherMixin from '@/mixins/watcher';
import { MODALS } from '@/constants';

const { mapActions } = createNamespacedHelpers('weatherWatcher');

export default {
  mixins: [
    modalMixin,
    watcherMixin,
  ],
  created() {
    this.fetchWeatherWatchersList({ filter: { _id: '5D0A5D2A4CAF0003' } });
  },
  methods: {
    ...mapActions({
      fetchWeatherWatchersList: 'fetchList',
    }),
    click() {
      this.showModal({
        name: MODALS.weatherWatcher,
        config: {
          watcherId: '5D0A5D2A4CAF0003',
        },
      });
    },
  },
};
</script>
