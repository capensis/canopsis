<template>
  <div :class="{ 'c-date-interval--column': column }" class="c-date-interval">
    <c-date-picker-field
      :value="value.from"
      :label="$t('common.from')"
      :disabled="disabled"
      :allowed-dates="isAllowedFromDate"
      hide-details
      @update:value="updateField('from', $event)"
    >
      <template #append-inner>
        <v-icon size="18">
          calendar_today
        </v-icon>
      </template>
    </c-date-picker-field>
    <c-date-picker-field
      :value="value.to"
      :label="$t('common.to')"
      :disabled="disabled"
      :allowed-dates="isAllowedToDate"
      :content-class="{ 'ml-4': !column }"
      hide-details
      @update:value="updateField('to', $event)"
    >
      <template #append-inner>
        <v-icon size="18">
          calendar_today
        </v-icon>
      </template>
    </c-date-picker-field>
  </div>
</template>

<script>
import { formMixin } from '@/mixins/form';

export default {
  emits: ['update:value'],
  mixins: [formMixin],
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

  &--column {
    flex-direction: column;
  }
}
</style>
