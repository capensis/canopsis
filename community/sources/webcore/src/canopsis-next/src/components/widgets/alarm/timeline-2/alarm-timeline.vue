<template>
  <v-layout column>
    <c-enabled-field
      :value="query.group"
      :label="$t('alarm.timeline.groupItems')"
      @input="updateGroup"
    />
    <alarm-timeline-days :days="days" />
    <c-pagination
      :total="meta.total_count"
      :limit="meta.per_page"
      :page="meta.page"
      @input="updatePage"
    />
  </v-layout>
</template>

<script>
import { groupAlarmSteps } from '@/helpers/entities/alarm/step/list';

import AlarmTimelineDays from './alarm-timeline-days.vue';

export default {
  components: {
    AlarmTimelineDays,
  },
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
  },
  computed: {
    days() {
      return groupAlarmSteps(this.steps);
    },
  },
  methods: {
    updateGroup(group) {
      this.$emit('update:query', {
        ...this.query,
        group,
      });
    },

    updatePage(page) {
      this.$emit('update:query', {
        ...this.query,
        page,
      });
    },
  },
};
</script>
