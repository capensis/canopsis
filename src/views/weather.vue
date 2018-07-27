<template lang="pug">
  div
    v-btn(@click="click") Open watcher modal
</template>

<script>
import modalMixin from '@/mixins/modal/modal';

export default {
  mixins: [
    modalMixin,
  ],
  methods: {
    click() {
      this.showModal({
        name: 'watcher-data',
      });
    },
  },
};
</script>
