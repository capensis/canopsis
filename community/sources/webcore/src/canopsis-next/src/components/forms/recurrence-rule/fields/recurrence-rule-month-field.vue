<template>
  <c-layout class="recurrence-rule-month-field" column>
    <div class="text-subtitle-2 mb-1">
      {{ t('recurrenceRule.bymonth') }}
    </div>

    <v-chip-group
      :model-value="modelValue"
      color="primary"
      multiple
      column
      @update:model-value="emit('update:modelValue', $event)"
    >
      <v-chip v-for="month in months" :key="month.value" :value="month.value">
        {{ month.text }}
      </v-chip>
    </v-chip-group>
  </c-layout>
</template>

<script setup>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

defineProps({
  modelValue: {
    type: Array,
    default: () => [],
  },
});

const emit = defineEmits(['update:modelValue']);

const { t } = useI18n();

const months = computed(() => [
  { text: t('common.months.january'), value: 1 },
  { text: t('common.months.february'), value: 2 },
  { text: t('common.months.march'), value: 3 },
  { text: t('common.months.april'), value: 4 },
  { text: t('common.months.may'), value: 5 },
  { text: t('common.months.june'), value: 6 },
  { text: t('common.months.july'), value: 7 },
  { text: t('common.months.august'), value: 8 },
  { text: t('common.months.september'), value: 9 },
  { text: t('common.months.october'), value: 10 },
  { text: t('common.months.november'), value: 11 },
  { text: t('common.months.december'), value: 12 },
]);
</script>
