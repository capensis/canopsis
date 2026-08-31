<template>
  <c-layout align-center>
    <div class="recurrence-rule-interval-field__input">
      <c-number-field
        :value="modelValue.interval"
        :label="t('common.every')"
        :min="1"
        name="interval"
        @update:value="updateInterval"
      />
    </div>

    <div class="pl-2 text-medium-emphasis">
      {{ intervalTimeString }}
    </div>
  </c-layout>
</template>

<script setup>
import { RRule } from 'rrule';
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
});

const emit = defineEmits(['update:modelValue']);

const { te, t } = useI18n();

const intervalTimeString = computed(() => {
  const timeMessageKey = {
    [RRule.HOURLY]: 'common.times.hour',
    [RRule.DAILY]: 'common.times.day',
    [RRule.WEEKLY]: 'common.times.week',
    [RRule.MONTHLY]: 'common.times.month',
    [RRule.YEARLY]: 'common.times.year',
  }[props.modelValue.freq];

  if (!timeMessageKey || !te(timeMessageKey)) {
    return '';
  }

  return t(timeMessageKey, props.modelValue.interval || 1);
});

const updateInterval = value =>
  emit('update:modelValue', {
    ...props.modelValue,
    interval: value,
  });
</script>

<style lang="scss" scoped>
.recurrence-rule-interval-field__input {
  flex: 0 0 50%;
}
</style>
