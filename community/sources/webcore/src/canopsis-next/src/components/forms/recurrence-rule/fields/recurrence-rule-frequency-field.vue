<template>
  <c-select
    :model-value="modelValue"
    :items="frequencies"
    :label="t('recurrenceRule.freq')"
    :autofocus="autofocus"
    item-title="text"
    item-value="value"
    @update:model-value="emit('update:modelValue', $event)"
  />
</template>

<script setup>
import { RRule } from 'rrule';
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

defineProps({
  modelValue: {
    type: Number,
    default: null,
  },
  autofocus: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update:modelValue']);

const { t } = useI18n();

const frequencies = computed(() => [
  { text: t('recurrenceRule.none'), value: null },
  { text: t('common.timeFrequencies.hourly'), value: RRule.HOURLY },
  { text: t('common.timeFrequencies.daily'), value: RRule.DAILY },
  { text: t('common.timeFrequencies.weekly'), value: RRule.WEEKLY },
  { text: t('common.timeFrequencies.monthly'), value: RRule.MONTHLY },
  { text: t('common.timeFrequencies.yearly'), value: RRule.YEARLY },
]);
</script>
