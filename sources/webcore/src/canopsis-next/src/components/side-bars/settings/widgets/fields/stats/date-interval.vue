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
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  methods: {
    showEditModal() {
      this.showModal({
        name: MODALS.statsDateInterval,
        config: {
          interval: this.value,
          action: date => this.$emit('input', date),
        },
      });
    },
  },
};
</script>
