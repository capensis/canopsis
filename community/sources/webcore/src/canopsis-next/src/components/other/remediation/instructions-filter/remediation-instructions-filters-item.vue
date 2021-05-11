<template lang="pug">
  div.instruction-filter
    v-chip.white--text(
      v-on="chipListeners",
      :color="chipColor",
      :close-icon="chipCloseIcon",
      :close="closable",
      label
    )
      span.instruction-filter__text
        v-icon(color="white", small) assignment
        v-icon.pl-1(v-if="filter.locked", color="white", small) lock
        strong.pl-2 {{ typeMessage }}
        span.pl-1(v-if="!filter.all") {{ instructionsMessage }}
</template>

<script>
import { MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

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
    closable: {
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

      return 'cancel';
    },

    anotherFilters() {
      return this.filters.filter(item => item._id !== this.filter._id);
    },

    typeMessage() {
      const getMessage = key => this.$t(`remediationInstructionsFilters.chip.${key}`);

      const { filter } = this;
      const all = filter.all || (filter.manual && filter.auto);

      return `${getMessage(filter.with ? 'with' : 'without')}${all ? ` ${getMessage('all')}` : ':'}`;
    },

    instructionsMessage() {
      return this.filter.instructions.map(({ name }) => name).join(', ');
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
          filters: this.anotherFilters,
          action: newFilter => this.updateModel(newFilter),
        },
      });
    },
  },
};
</script>

<style lang="scss">
.instruction-filter {
  & /deep/ .v-chip .v-chip__content {
    min-height: 32px;
    height: auto;
  }

  &__text {
    word-break: break-word;
    white-space: pre-line;
  }

  .v-chip__custom-close {
    font-size: 20px;
  }
}
</style>
