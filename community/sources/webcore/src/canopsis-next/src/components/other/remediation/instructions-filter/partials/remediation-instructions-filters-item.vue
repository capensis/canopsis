<template>
  <div class="instruction-filter">
    <v-chip
      class="white--text"
      v-on="chipListeners"
      :color="chipColor"
      :close-icon="chipCloseIcon"
      :close="closable"
      label
    >
      <span class="instruction-filter__text">
        <v-icon
          color="white"
          small
        >
          assignment
        </v-icon>
        <v-icon
          class="pl-1"
          v-if="filter.locked"
          color="white"
          small
        >
          lock
        </v-icon>
        <strong class="pl-2 text-uppercase">{{ conditionTypeMessage }}</strong>
        <span
          class="pl-1"
          v-if="!isAll"
        >
          {{ typesAndInstructionsMessage }}
        </span>
        <strong
          class="pl-1 text-uppercase"
          v-if="hasRunning"
        >
          {{ $t('remediation.instructionsFilter.inProgress') }}
        </strong>
      </span>
    </v-chip>
  </div>
</template>

<script>
import { isBoolean } from 'lodash';

import { MODALS, REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

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
    isAll() {
      return this.filter.all || (this.filter.manual && this.filter.auto);
    },

    chipListeners() {
      const listeners = { 'click:close': this.close };

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

    instructionsNames() {
      return this.filter.instructions?.map(({ name }) => name) ?? [];
    },

    hasRunning() {
      return isBoolean(this.filter.running);
    },

    conditionalTypeMessagePrefix() {
      if (!this.hasRunning) {
        return '';
      }

      const message = this.filter.running ? this.$t('common.show') : this.$t('common.hide');

      return `${message} `;
    },

    conditionTypeMessage() {
      const allMessage = this.isAll ? ` ${this.$t('remediation.instructionsFilter.chip.all')}` : ':';
      const conditionMessage = this.$t(`remediation.instructionsFilter.chip.${this.filter.with ? 'with' : 'without'}`);

      return `${this.conditionalTypeMessagePrefix}${conditionMessage}${allMessage}`;
    },

    typesAndInstructionsMessage() {
      const types = [];

      if (this.filter.manual) {
        types.push(this.$t(`remediation.instruction.types.${REMEDIATION_INSTRUCTION_TYPES.manual}`));
      }

      if (this.filter.auto) {
        types.push(this.$t(`remediation.instruction.types.${REMEDIATION_INSTRUCTION_TYPES.auto}`));
      }

      return [...types, ...this.instructionsNames].join(', ');
    },
  },
  methods: {
    close() {
      if (this.filter.locked) {
        return this.updateField('disabled', !this.filter.disabled);
      }

      return this.$emit('remove');
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
  & .v-chip .v-chip__content {
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
