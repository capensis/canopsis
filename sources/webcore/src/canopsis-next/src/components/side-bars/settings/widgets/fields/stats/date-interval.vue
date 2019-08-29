<template lang="pug">
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      div.subheading Date interval
      div
        v-btn.primary(
          small,
          @click="showEditModal"
        ) {{ $t('common.edit') }}
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
    hiddenFields: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    showEditModal() {
      this.showModal({
        name: MODALS.statsDateInterval,
        config: {
          interval: this.value,
          hiddenFields: this.hiddenFields,
          action: date => this.$emit('input', date),
        },
      });
    },
  },
};
</script>
