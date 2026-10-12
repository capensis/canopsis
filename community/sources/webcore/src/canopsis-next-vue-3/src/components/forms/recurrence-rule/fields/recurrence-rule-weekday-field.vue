<template>
  <c-layout v-if="chips" class="recurrence-rule-weekday-field" column>
    <div class="text-subtitle-2 mb-1">
      {{ t('recurrenceRule.byweekday') }}
    </div>

    <v-chip-group
      :model-value="modelValue"
      color="primary"
      multiple
      @update:model-value="emit('update:modelValue', $event)"
    >
      <v-chip v-for="weekDay in weekDays" :key="weekDay.value" :value="weekDay.value">
        {{ weekDay.text }}
      </v-chip>
    </v-chip-group>
  </c-layout>

  <c-select
    v-else
    :model-value="modelValue"
    :items="weekDays"
    :label="t('recurrenceRule.wkst')"
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
    type: [Array, Number, String],
    default: () => [],
  },
  chips: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update:modelValue']);

const { t } = useI18n();

const weekDays = computed(() => [
  { text: t('common.weekDays.monday'), value: RRule.MO.weekday },
  { text: t('common.weekDays.tuesday'), value: RRule.TU.weekday },
  { text: t('common.weekDays.wednesday'), value: RRule.WE.weekday },
  { text: t('common.weekDays.thursday'), value: RRule.TH.weekday },
  { text: t('common.weekDays.friday'), value: RRule.FR.weekday },
  { text: t('common.weekDays.saturday'), value: RRule.SA.weekday },
  { text: t('common.weekDays.sunday'), value: RRule.SU.weekday },
]);
</script>
