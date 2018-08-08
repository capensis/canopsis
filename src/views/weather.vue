<template lang="pug">
  div
    v-btn(@click="click") Open watcher modal
</template>

<script>
import modalMixin from '@/mixins/modal/modal';
import watcherMixin from '@/mixins/watcher';
import { MODALS } from '@/constants';

export default {
  mixins: [
    modalMixin,
    watcherMixin,
  ],
  created() {
    this.fetchWatchersList({ filter: { _id: 'AE05C0454680067A' } });
  },
  methods: {
    click() {
      this.showModal({
        name: MODALS.weatherWatcher,
        config: {
          watcherId: 'AE05C0454680067A',
        },
      });
    },
  },
};
</script>
