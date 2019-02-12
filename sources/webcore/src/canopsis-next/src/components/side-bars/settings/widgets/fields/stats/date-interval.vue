<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") Date interval
    v-container
      v-btn(@click="showEditModal") Edit
</template>

<script>
import modalMixin from '@/mixins/modal';
import { MODALS } from '@/constants';

export default {
  mixins: [modalMixin],
  methods: {
    showEditModal() {
      this.showModal({
        name: MODALS.statsDateInterval,
      });
    },
  },
};
</script>
