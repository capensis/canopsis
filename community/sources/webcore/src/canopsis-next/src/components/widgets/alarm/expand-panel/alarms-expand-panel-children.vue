<template lang="pug">
  div
    alarms-list-table(
      :widget="widget",
      :alarms="alarms",
      :total-items="meta.total_count",
      :pagination.sync="pagination",
      :editing="editing",
      :columns="columns",
      :loading="pending",
      :parent-alarm="alarm",
      expandable,
      hide-groups,
      ref="alarmsTable"
    )
    c-table-pagination(
      :total-items="meta.total_count",
      :rows-per-page="query.limit",
      :page="query.page",
      @update:page="updateQueryPage",
      @update:rows-per-page="updateRecordsPerPage"
    )
</template>

<script>
import { isEqual, pick } from 'lodash';

import { DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS, ALARM_ENTITY_FIELDS, SORT_ORDERS } from '@/constants';

import { defaultColumnsToColumns } from '@/helpers/entities';

/**
 * Group-alarm-list component
 *
 * @module alarm
 *
 */
export default {
  props: {
    children: {
      type: Object,
      required: true,
    },
    alarm: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    query: {
      type: Object,
      required: true,
    },
    editing: {
      type: Boolean,
      default: false,
    },
    pending: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    alarms() {
      return this.children?.data ?? [];
    },

    meta() {
      return this.children?.meta ?? {};
    },

    columns() {
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

    pagination: { // TODO: move to mixin
      get() {
        const { sortDir, sortKey: sortBy = null, multiSortBy = [] } = this.query;
        const descending = sortDir === SORT_ORDERS.desc;

        return { sortBy, descending, multiSortBy };
      },

      set(value) {
        const paginationKeys = ['sortBy', 'descending', 'multiSortBy'];
        const newPagination = pick(value, paginationKeys);
        const oldPagination = pick(this.pagination, paginationKeys);

        if (isEqual(newPagination, oldPagination)) {
          return;
        }

        const {
          sortBy = null,
          descending = false,
          multiSortBy = [],
        } = newPagination;

        const newQuery = {
          sortKey: sortBy,
          sortDir: descending ? SORT_ORDERS.desc : SORT_ORDERS.asc,
          multiSortBy,
        };

        this.$emit('update:query', {
          ...this.query,
          ...newQuery,
        });
      },
    },
  },
  methods: {
    updateRecordsPerPage(limit) {
      this.$emit('update:query', {
        ...this.query,

        limit,
      });
    },

    updateQueryPage(page) {
      this.$emit('update:query', {
        ...this.query,

        page,
      });
    },
  },
};
</script>
