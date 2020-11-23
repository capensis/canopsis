<template lang="pug">
  div
    v-chip.white--text(
      v-on="chipListeners",
      :color="chipColor",
      :close-icon="chipCloseIcon",
      close,
      label
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
    editable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    chipListeners() {
      const listeners = { input: this.close };

      if (this.editable) {
        listeners.click = this.showEditFilterModal;
      }

      return listeners;
    },

    chipColor() {
      return this.filter.disabled ? 'grey' : 'primary';
    },

    chipCloseIcon() {
      if (this.filter.locked) {
        return this.filter.disabled ? 'check_box_outline_blank' : 'check_box';
      }

      return '$vuetify.icons.delete';
    },

    anotherFilters() {
      return this.filters.filter(item => item._id !== this.filter._id);
    },

    typeMessage() {
      const getMessage = key => this.$t(`remediationInstructionsFilters.chip.${key}`);

      const { filter } = this;

      return `${filter.with ? getMessage('with') : getMessage('without')}${filter.all ? ` ${getMessage('all')}` : ':'}`;
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
