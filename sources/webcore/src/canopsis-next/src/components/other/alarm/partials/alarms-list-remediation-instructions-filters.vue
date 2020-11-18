<template lang="pug">
  div
    v-chip.primary.white--text(
      v-for="filter in filters",
      :key="filter._id",
      close,
      label,
      @click="showEditFilterModal(filter)",
      @input="removeFilter(filter)"
    )
      span
        strong {{ filter | conditionMessage }}
        span.pl-1(v-if="!filter.all") {{ filter.instructions.join(', ') }}
    v-tooltip(bottom)
      v-btn(
        slot="activator",
        icon,
        small,
        @click="showCreateFilterModal"
      )
        v-icon(:color="filters.length ? 'primary' : 'black'") adjust
      span {{ $t('instructionsFilter.button') }}
</template>

<script>
import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

export default {
  filters: {
    conditionMessage({ condition, all }) {
      return `${condition === 0 ? 'WITH' : 'WITHOUT'}${all ? ' ALL' : ':'}`;
    },
  },
  props: {
    filters: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    removeFilter(filter) {
      this.$emit('input', this.filters.filter(item => item._id !== filter._id));
    },
    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.remediationInstructionsFilterEditor,
        config: {
          action: newFilter => this.$emit('input', [...this.filters, { _id: uid(), ...newFilter }]),
        },
      });
    },
    showEditFilterModal(filter) {
      this.$modals.show({
        name: MODALS.remediationInstructionsFilterEditor,
        config: {
          filter,
          action: newFilter =>
            this.$emit('input', this.filters.map(item => (item._id === filter._id ? newFilter : item))),
        },
      });
    },
  },
};
</script>
