<template>
  <alarms-list-table
    :widget="widget"
    :alarms="alarms"
    :total-items="meta.total_count"
    :options.sync="options"
    :columns="columns"
    :loading="loading"
    :parent-alarm="parentAlarm"
    :hide-children="hideChildren"
    :sticky-header="stickyHeader"
    :refresh-alarms-list="refreshAlarmsList"
    :selectable="selectable"
    :expandable="expandable"
    :hide-actions="hideActions"
    :hide-pagination="hidePagination"
    @update:page="updatePage"
    @update:items-per-page="updateItemsPerPage"
  />
</template>

<script>
import { PAGINATION_LIMIT } from '@/config';

import { convertDataTableOptionsToQuery } from '@/helpers/entities/shared/query';
import { getPageForNewItemsPerPage } from '@/helpers/pagination';

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
    loading: {
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
    options: {
      get() {
        const { page = 1, itemsPerPage = PAGINATION_LIMIT, sortBy = [], sortDesc = [] } = this.query;

        return { page, itemsPerPage, sortBy, sortDesc };
      },

      set(newOptions) {
        const convertedOptions = convertDataTableOptionsToQuery(newOptions, this.options);

        if (convertedOptions === this.options) {
          return;
        }

        this.$emit('update:query', {
          ...this.query,
          ...convertedOptions,
        });
      },
    },
  },
  methods: {
    updateItemsPerPage(itemsPerPage) {
      this.$emit('update:query', {
        ...this.query,

        itemsPerPage,
        page: getPageForNewItemsPerPage(itemsPerPage, this.query.itemsPerPage, this.query.page),
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
