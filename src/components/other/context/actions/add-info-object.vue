<template>
  <v-btn @click="showForm">Add info</v-btn>
</template>

<script>
import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

export default {
  mixins: [
    modalMixin,
  ],
  props: {
    entity: {
      type: Object,
      required: true,
    },
  },
  methods: {
    showForm() {
      this.showModal({
        name: MODALS.infoObject,
        config: {
          entity: this.entity,
        },
      });
    },
  },
};
</script>
