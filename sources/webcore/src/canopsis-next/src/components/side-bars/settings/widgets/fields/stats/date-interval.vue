<template lang="pug">
  v-container.pa-3(fluid, data-test="dateInterval")
    v-layout(align-center, justify-space-between)
      div.subheading Date interval
      div
        v-btn.primary(
          data-test="editButton",
          small,
          @click="showEditModal"
        ) {{ $t('common.edit') }}
</template>

<script>
import { MODALS } from '@/constants';

export default {
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
      this.$modals.show({
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
