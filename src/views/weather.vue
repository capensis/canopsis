<template lang="pug">
  div
    v-btn(@click="click") Open watcher modal
</template>

<script>
import modalMixin from '@/mixins/modal/modal';
import watcherMixin from '@/mixins/watcher';

export default {
  mixins: [
    modalMixin,
    watcherMixin,
  ],
  created() {
    this.fetchWatchersList();
  },
  methods: {
    click() {
      this.showModal({
        name: 'watcher-data',
        config: {
          watcherId: 'test-watcher',
        },
      });
    },
  },
};
</script>
