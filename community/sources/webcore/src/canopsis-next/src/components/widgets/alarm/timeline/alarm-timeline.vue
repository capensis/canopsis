<template>
  <v-layout column>
    <v-layout class="gap-5">
      <c-enabled-field
        :value="query.group"
        :label="$t('alarm.timeline.groupItems')"
        @input="updateGroup"
      />
      <c-enabled-field
        :value="isCommentType"
        :label="$t('alarm.timeline.onlyComments')"
        @input="updateOnlyComments"
      />
    </v-layout>
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

import { ALARM_STEPS_TYPES } from '@/constants';

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
    const isCommentType = computed(() => props.query.type === ALARM_STEPS_TYPES.comment);

    const updateGroup = group => emit('update:query', { ...props.query, group, page: 1 });
    const updateOnlyComments = (onlyComments) => {
      const { type, ...newQuery } = props.query;

      newQuery.page = 1;

      if (onlyComments) {
        newQuery.type = ALARM_STEPS_TYPES.comment;
      }

      emit('update:query', newQuery);
    };

    const updatePage = page => emit('update:query', { ...props.query, page });

    return {
      days,
      isCommentType,

      updateGroup,
      updateOnlyComments,
      updatePage,
    };
  },
};
</script>
