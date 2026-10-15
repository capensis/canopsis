<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ title }}</span>
    </template>

    <template #text>
      <pbehaviors-calendar-component :entity-id="config.entityId" />
    </template>
  </c-modal-wrapper>
</template>

<script setup>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import PbehaviorsCalendarComponent from '@/components/other/pbehavior/calendar/pbehaviors-calendar.vue';

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();

const config = computed(() => props.modal.config ?? {});
const title = computed(() => config.value.title ?? t('modals.pbehaviorsCalendar.title'));
</script>
