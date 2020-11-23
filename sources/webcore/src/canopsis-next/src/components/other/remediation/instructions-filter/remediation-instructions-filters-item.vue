<template lang="pug">
  div
    v-chip.white--text(
      :color="color",
      :close-icon="closeIcon",
      close,
      label,
      @click="showEditFilterModal",
      @input="close"
    )
      span
        v-icon(color="white", small) assignment
        v-icon.pl-1(v-if="filter.locked", color="white", small) lock
        strong.pl-2 {{ typeMessage }}
        span.pl-1(v-if="!filter.all") {{ instructionsMessage }}
</template>

<script>
import { MODALS } from '@/constants';

import formMixin from '@/mixins/form';

export default {
  mixins: [formMixin],
  model: {
    prop: 'filter',
    event: 'input',
  },
  props: {
    filter: {
      type: Object,
      default: () => ({}),
    },
    filters: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    color() {
      return this.filter.disabled ? 'grey' : 'primary';
    },
    closeIcon() {
      if (this.filter.locked) {
        return this.filter.disabled ? 'check_box_outline_blank' : 'check_box';
      }

      return '$vuetify.icons.delete';
    },
    anotherFilters() {
      return this.filters.filter(item => item._id !== this.filter._id);
    },
    typeMessage() {
      const { filter } = this;

      return `${filter.with ? 'WITH' : 'WITHOUT'}${filter.all ? ' ALL' : ':'}`; // TODO: add i18n
    },
    instructionsMessage() {
      return this.filter.instructions.join(', ');
    },
  },
  methods: {
    close() {
      if (this.filter.locked) {
        return this.updateField('disabled', !this.filter.disabled);
      }

      return this.$emit('remove', this.filter);
    },
    showEditFilterModal() {
      this.$modals.show({
        name: MODALS.createRemediationInstructionsFilter,
        config: {
          filter: this.filter,
          anotherFilters: this.anotherFilters,
          action: newFilter => this.updateModel(newFilter),
        },
      });
    },
  },
};
</script>

<style lang="scss">
.v-chip__custom-close {
  font-size: 20px;
}
</style>
