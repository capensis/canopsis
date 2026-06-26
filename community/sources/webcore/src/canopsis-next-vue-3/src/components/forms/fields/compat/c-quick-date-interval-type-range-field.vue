<template>
  <c-layout gap="3">
    <c-quick-date-interval-type-field
      :value="value.from"
      :ranges="fromRanges"
      :label="$t('common.from')"
      @update:value="updateField('from', $event)"
    />
    <c-quick-date-interval-type-field
      :value="value.to"
      :ranges="toRanges"
      :label="$t('common.to')"
      @update:value="updateField('to', $event)"
    />
  </c-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';

export default {
  emits: ['update:value'],
  mixins: [formMixin],
  props: {
    value: {
      type: Object,
      required: true,
    },
    fromRanges: {
      type: Array,
      default: () => [],
    },
    toRanges: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
