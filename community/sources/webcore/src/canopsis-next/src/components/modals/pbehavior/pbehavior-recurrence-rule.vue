<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ config.pbehavior?.name }}</span>
    </template>

    <template #text>
      <pbehavior-recurrence-rule-calendar :pbehavior="config.pbehavior" />
    </template>
  </c-modal-wrapper>
</template>

<script setup>
import { computed } from 'vue';

import PbehaviorRecurrenceRuleCalendar from '@/components/other/pbehavior/calendar/pbehavior-recurrence-rule-calendar.vue';

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const config = computed(() => props.modal.config ?? {});
</script>
