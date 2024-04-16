<template>
  <div :class="['c-date-time-interval-field', 'col-gap-3', { 'c-date-time-interval-field--column': column }]">
    <date-time-picker-field
      v-validate="rules"
      v-field="value.from"
      :label="$t('common.from')"
      :disabled="disabled"
      :name="fromFieldName"
      :hide-details="hideDetails"
      :allowed-dates="isAllowedFromDate"
      :round-hours="roundHours"
    />
    <date-time-picker-field
      v-validate="rules"
      v-field="value.to"
      :label="$t('common.to')"
      :disabled="disabled"
      :name="toFieldName"
      :hide-details="hideDetails"
      :allowed-dates="isAllowedToDate"
      :round-hours="roundHours"
    />
  </div>
</template>

<script>
import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DateTimePickerField,
  },
  model: {
    event: 'input',
    prop: 'value',
  },
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
