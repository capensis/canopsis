<template>
  <v-layout column>
    <c-enabled-field
      :value="query.group"
      :label="$t('alarm.timeline.groupItems')"
      @input="updateGroup"
    />
    <alarm-timeline-days :days="days" :is-html-enabled="isHtmlEnabled" />
    <c-pagination
      :total="meta.total_count"
      :limit="meta.per_page"
      :page="meta.page"
      @input="updatePage"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { groupAlarmSteps } from '@/helpers/entities/alarm/step/list';

import AlarmTimelineDays from './alarm-timeline-days.vue';

export default {
  components: { AlarmTimelineDays },
  props: {
    steps: {
      type: Array,
      default: () => [],
    },
    meta: {
      type: Object,
      default: () => ({}),
    },
    query: {
      type: Object,
      default: () => ({}),
    },
    isHtmlEnabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const days = computed(() => groupAlarmSteps(props.steps));

    const updateGroup = group => emit('update:query', { ...props.query, group, page: 1 });
    const updatePage = page => emit('update:query', { ...props.query, page });

    return {
      days,

      updateGroup,
      updatePage,
    };
  },
};
</script>
