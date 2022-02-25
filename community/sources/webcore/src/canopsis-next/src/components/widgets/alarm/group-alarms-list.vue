<template lang="pug">
  div
    alarms-list-table(
      :widget="widget",
      :alarms="displayedAlarms",
      :total-items="alarmsMeta.total_count",
      :pagination.sync="pagination",
      :editing="editing",
      :has-columns="hasGroupColumns",
      :columns="groupColumns",
      :parent-alarm="alarm",
      expandable,
      hide-groups,
      ref="alarmsTable"
    )
    c-table-pagination(
      :total-items="alarmsMeta.total_count",
      :rows-per-page="query.limit",
      :page="query.page",
      @update:page="updateQueryPage",
      @update:rows-per-page="updateRecordsPerPage"
    )
</template>

<script>
import { get, orderBy } from 'lodash';

import { DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS, ALARM_ENTITY_FIELDS } from '@/constants';

import { defaultColumnsToColumns } from '@/helpers/entities';
import { convertWidgetToQuery } from '@/helpers/query';

import { queryWidgetMixin } from '@/mixins/widget/query';

/**
 * Group-alarm-list component
 *
 * @module alarm
 *
 */
export default {
  mixins: [queryWidgetMixin],
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    editing: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    alarms() {
      return get(this.alarm, 'consequences.data') || get(this.alarm, 'causes.data', []);
    },

    alarmsMeta() {
      return {
        total_count: this.alarms.length,
      };
    },

    displayedAlarms() {
      const {
        page,
        limit,
        multiSortBy = [],
      } = this.query;

      let { alarms } = this;

      if (multiSortBy.length) {
        alarms = orderBy(
          alarms,
          multiSortBy.map(({ sortBy }) => sortBy),
          multiSortBy.map(({ descending }) => (descending ? 'desc' : 'asc')),
        );
      }

      return alarms.slice((page - 1) * limit, page * limit);
    },

    hasGroupColumns() {
      return this.groupColumns.length > 0;
    },

    groupColumns() {
      if (this.widget.parameters.widgetGroupColumns) {
        return this.widget.parameters.widgetGroupColumns.map(({ value, label, ...column }) => ({
          ...column,
          value,
          text: label,
          sortable: value !== ALARM_ENTITY_FIELDS.extraDetails,
        }));
      }

      return defaultColumnsToColumns(DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS);
    },
  },
  mounted() {
    this.query = convertWidgetToQuery(this.widget);
  },
};
</script>
