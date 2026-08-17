<template>
  <div
    :class="[
      'c-date-time-interval-field',
      'col-gap-3',
      { 'c-date-time-interval-field--column': column },
    ]"
  >
    <date-time-picker-field
      :value="value.from"
      :label="$t('common.from')"
      :disabled="disabled"
      :name="fromFieldName"
      :hide-details="hideDetails"
      :allowed-dates="isAllowedFromDate"
      :round-hours="roundHours"
      :rules="rules"
      :required="required"
      @update:value="updateField('from', $event)"
    >
      <template v-if="withIcon" #append-inner>
        <v-icon size="18">
          calendar_today
        </v-icon>
      </template>
    </date-time-picker-field>
    <date-time-picker-field
      :value="value.to"
      :label="$t('common.to')"
      :disabled="disabled"
      :name="toFieldName"
      :hide-details="hideDetails"
      :allowed-dates="isAllowedToDate"
      :round-hours="roundHours"
      :rules="rules"
      :required="required"
      @update:value="updateField('to', $event)"
    >
      <template v-if="withIcon" #append-inner>
        <v-icon size="18">
          calendar_today
        </v-icon>
      </template>
    </date-time-picker-field>
  </div>
</template>

<script>
import { formMixin } from '@/mixins/form';

import DateTimePickerField from './date-time-picker-field.vue';

export default {
  components: {
    DateTimePickerField,
  },
  emits: ['update:value'],
  mixins: [formMixin],
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'interval',
    },
    disabled: {
      type: Boolean,
      required: false,
    },
    hideDetails: {
      type: Boolean,
      required: false,
    },
    required: {
      type: Boolean,
      required: false,
    },
    isAllowedFromDate: {
      type: Function,
      required: false,
    },
    isAllowedToDate: {
      type: Function,
      required: false,
    },
    column: {
      type: Boolean,
      default: false,
    },
    roundHours: {
      type: Boolean,
      default: false,
    },
    withIcon: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

    fromFieldName() {
      return `${this.name}.from`;
    },

    toFieldName() {
      return `${this.name}.to`;
    },
  },
};
</script>

<style lang="scss">
.c-date-time-interval-field {
  display: inline-flex;

  &--column {
    flex-direction: column;
  }
}
</style>
