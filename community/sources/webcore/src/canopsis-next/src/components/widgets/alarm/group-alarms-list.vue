<template lang="pug">
  div
    alarms-list-table(
      :widget="widget",
      :alarms="displayedAlarms",
      :total-items="alarmsMeta.total_count",
      :pagination.sync="vDataTablePagination",
      :editing="editing",
      :has-columns="hasGroupColumns",
      :columns="groupColumns",
      :parent-alarm="alarm",
      expandable,
      hideGroups,
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
import { orderBy } from 'lodash';

import { DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS, ALARM_ENTITY_FIELDS } from '@/constants';

import { defaultColumnsToColumns } from '@/helpers/entities';

import { widgetColumnsAlarmMixin } from '@/mixins/widget/columns';
import widgetGroupFetchQueryMixin from '@/mixins/widget/group-fetch-query';
import widgetExpandPanelAlarm from '@/mixins/widget/expand-panel/alarm/expand-panel';

/**
 * Group-alarm-list component
 *
 * @module alarm
 *
 */
export default {
  mixins: [
    widgetGroupFetchQueryMixin,
    widgetColumnsAlarmMixin,
    widgetExpandPanelAlarm,
  ],
  props: {
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

    displayedAlarms() {
      const {
        page,
        limit,
        multiSortBy = [],
      } = this.query;

      let { alarms } = this;
      const filteredMultiSortBy = multiSortBy
        .filter(({ sortBy }) => this.groupColumns.some(({ value }) => value.endsWith(sortBy)));

      if (filteredMultiSortBy.length) {
        alarms = orderBy(
          alarms,
          multiSortBy.map(({ sortBy }) => sortBy),
          multiSortBy.map(({ descending }) => (descending ? 'desc' : 'asc')),
        );
      }

      return alarms.slice((page - 1) * limit, page * limit);
    },
  },
};
</script>
