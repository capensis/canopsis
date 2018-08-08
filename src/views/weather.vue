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
    this.fetchWeatherWatchersList({ filter: { _id: 'BD41E18E45870569' } });
  },
  methods: {
    ...mapActions({
      fetchWeatherWatchersList: 'fetchList',
    }),
    click() {
      this.showModal({
        name: MODALS.weatherWatcher,
        config: {
          watcherId: 'BD41E18E45870569',
        },
      });
    },
  },
};
</script>
