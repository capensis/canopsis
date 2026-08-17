<template>
  <c-layout class="position-relative" column>
    <c-progress-overlay :pending="pending" />
    <c-advanced-data-table
      :items="ganttIntervals"
      :headers="headers"
      :total-items="meta.total_count"
      :options="query"
      advanced-pagination
      hide-actions
      @update:options="query = $event"
    >
      <template #from="{ item }">
        {{ dateToString(item.from) }}
      </template>
      <template #to="{ item }">
        {{ dateToString(item.to) }}
      </template>
    </c-advanced-data-table>
  </c-layout>
</template>

<script>
import { mapActions } from 'pinia';

import { PAGINATION_LIMIT } from '@/config';
import { DATETIME_FORMATS } from '@/constants';

import { useTestSuiteStore } from '@/stores/modules/test-suite';

import { convertDateToString } from '@/helpers/date/date';

export default {
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },

  data() {
    return {
      pending: false,
      meta: {},
      ganttIntervals: [],
      query: {
        itemsPerPage: PAGINATION_LIMIT,
        page: 1,
      },
    };
  },

  computed: {
    headers() {
      return [
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.status'), value: 'status' },
        { text: this.$t('common.from'), value: 'from' },
        { text: this.$t('common.to'), value: 'to' },
        { text: this.$t('common.duration'), value: 'time' },
      ];
    },
  },

  watch: {
    query: {
      deep: true,
      immediate: true,
      handler() {
        this.fetchList();
      },
    },
  },

  methods: {
    ...mapActions(useTestSuiteStore, ['fetchEntityGanttIntervalsWithoutStore']),

    dateToString(value) {
      return convertDateToString(value, DATETIME_FORMATS.long, '-');
    },

    async fetchList() {
      this.pending = true;

      try {
        const { page, itemsPerPage: limit } = this.query;
        const { data, meta } = await this.fetchEntityGanttIntervalsWithoutStore({
          id: this.alarm.entity._id,
          params: { page, limit },
        });

        this.ganttIntervals = data;
        this.meta = meta;
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
