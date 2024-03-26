<template>
  <div
    :class="{ 'c-date-interval--column': column }"
    class="c-date-interval"
  >
    <c-date-picker-field
      v-field="value.from"
      :label="$t('common.from')"
      :disabled="disabled"
      :allowed-dates="isAllowedFromDate"
      hide-details
    >
      <template #append="">
        <v-icon>calendar_today</v-icon>
      </template>
    </c-date-picker-field>
    <c-date-picker-field
      v-field="value.to"
      :label="$t('common.to')"
      :disabled="disabled"
      :allowed-dates="isAllowedToDate"
      :content-class="{ 'ml-4': !column }"
      hide-details
    >
      <template #append="">
        <v-icon>calendar_today</v-icon>
      </template>
    </c-date-picker-field>
  </div>
</template>

<script>
export default {
  model: {
    event: 'input',
    prop: 'value',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    isAllowedFromDate: {
      type: Function,
      required: false,
    },
    isAllowedToDate: {
      type: Function,
      required: false,
    },
    disabled: {
      type: Boolean,
      required: false,
    },
    column: {
      type: Boolean,
      required: false,
    },
  },
};
</script>

<style scoped lang="scss">
.c-date-interval {
  display: inline-flex;

  .v-input__append-inner {
    pointer-events: none;
  }

  &--column {
    flex-direction: column;
  }
}
</style>
