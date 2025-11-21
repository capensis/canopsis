<template>
  <alarms-list-table
    :widget="widget"
    :alarms="alarms"
    :total-items="meta.total_count"
    :options.sync="options"
    :columns="columns"
    :columns-settings="columnsSettings"
    :loading="loading"
    :parent-alarm="parentAlarm"
    :hide-children="hideChildren"
    :sticky-header="stickyHeader"
    :refresh-alarms-list="refreshAlarmsList"
    :selectable="selectable"
    :expandable="expandable"
    :hide-actions="hideActions"
    :hide-pagination="hidePagination"
    :resizable-column="resizableColumn"
    :draggable-column="draggableColumn"
    :cells-content-behavior="cellsContentBehavior"
    eager
    @update:page="updatePage"
    @update:items-per-page="updateItemsPerPage"
    @update:columns-settings="updateColumnsSettings"
  />
</template>

<script>
import { computed } from 'vue';

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
    columnsSettings: {
      type: Object,
      default: () => ({}),
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
  setup(props, { emit }) {
    const resizableColumn = computed(() => !!props.widget.parameters?.columns?.resizable);
    const draggableColumn = computed(() => !!props.widget.parameters?.columns?.draggable);
    const cellsContentBehavior = computed(() => props.widget.parameters?.columns?.cells_content_behavior);

    const options = computed({
      get() {
        const { page = 1, itemsPerPage = PAGINATION_LIMIT, sortBy = [], sortDesc = [] } = props.query;

        return { page, itemsPerPage, sortBy, sortDesc };
      },

      set(newOptions) {
        const convertedOptions = convertDataTableOptionsToQuery(newOptions, options.value);

        if (convertedOptions === options.value) {
          return;
        }

        emit('update:query', {
          ...props.query,
          ...convertedOptions,
        });
      },
    });

    /**
     * Updates items per page and recalculates the current page number
     *
     * @param {number} itemsPerPage - New items per page value
     */
    const updateItemsPerPage = itemsPerPage => emit('update:query', {
      ...props.query,

      itemsPerPage,
      page: getPageForNewItemsPerPage(itemsPerPage, props.query.itemsPerPage, props.query.page),
    });

    /**
     * Updates the current page number
     *
     * @param {number} page - New page number
     */
    const updatePage = page => emit('update:query', {
      ...props.query,

      page,
    });

    /**
     * Updates the columns settings
     *
     * @param {Object} columnsSettings - New columns settings
     */
    const updateColumnsSettings = columnsSettings => emit('update:columns-settings', columnsSettings);

    return {
      resizableColumn,
      draggableColumn,
      cellsContentBehavior,
      options,
      updateItemsPerPage,
      updatePage,
      updateColumnsSettings,
    };
  },
};
</script>
