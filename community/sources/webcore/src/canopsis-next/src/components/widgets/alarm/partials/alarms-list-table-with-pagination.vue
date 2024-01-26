<template lang="pug">
  alarms-list-table(
    :widget="widget",
    :alarms="alarms",
    :total-items="meta.total_count",
    :pagination.sync="pagination",
    :columns="columns",
    :editing="editing",
    :loading="loading",
    :parent-alarm="parentAlarm",
    :is-tour-enabled="isTourEnabled",
    :hide-children="hideChildren",
    :sticky-header="stickyHeader",
    :refresh-alarms-list="refreshAlarmsList",
    :selectable="selectable",
    :expandable="expandable",
    :hide-actions="hideActions",
    :hide-pagination="hidePagination",
    @update:page="updateQueryPage",
    @update:rows-per-page="updateRecordsPerPage"
  )
</template>

<script>
import { isEqual, pick } from 'lodash';

import { getPageForNewRecordsPerPage } from '@/helpers/pagination';

import { SORT_ORDERS } from '@/constants';

/**
 * Group-alarm-list component
 *
 * @module alarm
 *
 */
export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
    alarms: {
      type: Array,
      required: true,
    },
    columns: {
      type: Array,
      default: () => [],
    },
    meta: {
      type: Object,
      required: true,
    },
    query: {
      type: Object,
      required: true,
    },
    parentAlarm: {
      type: Object,
      default: null,
    },
    editing: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    isTourEnabled: {
      type: Boolean,
      default: false,
    },
    expandable: {
      type: Boolean,
      default: false,
    },
    selectable: {
      type: Boolean,
      default: false,
    },
    stickyHeader: {
      type: Boolean,
      default: false,
    },
    hideChildren: {
      type: Boolean,
      default: false,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
    hidePagination: {
      type: Boolean,
      default: false,
    },
    refreshAlarmsList: {
      type: Function,
      default: () => {},
    },
  },
  computed: {
    pagination: {
      get() {
        const { sortDir, page, limit, sortKey: sortBy = null, multiSortBy = [] } = this.query;
        const descending = sortDir === SORT_ORDERS.desc;

        return { sortBy, page, limit, descending, multiSortBy };
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
        page: getPageForNewRecordsPerPage(limit, this.query.limit, this.query.page),
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
