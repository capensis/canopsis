<template>
  <div class="alarms-list-table-with-pagination">
    <alarms-list-table
      :widget="widget"
      :alarms="alarms"
      :columns="columns"
      :columns-settings="columnsSettings"
      :meta="meta"
      :query="query"
      :parent-alarm="parentAlarm"
      :loading="loading"
      :expandable="expandable"
      :selectable="selectable"
      :sticky-header="stickyHeader"
      :sticky-horizontal-scroll="stickyHorizontalScroll"
      :hide-children="hideChildren"
      :hide-actions="hideActions"
      :densable="densable"
      :dense-value="dense"
      :widget-grid-parameters="widgetGridParameters"
      :default-keep-selected-after-action="defaultKeepSelectedAfterAction"
      :refresh-alarms-list="refreshAlarmsList"
      :selected-tags="selectedTags"
      :eager="eager"
      @update:query="updateQuery"
      @update:columns-settings="updateColumnsSettings"
      @update:dense="$emit('update:dense', $event)"
      @select:tag="$emit('select:tag', $event)"
      @remove:tag="$emit('remove:tag', $event)"
    />

    <c-table-pagination
      v-if="!hidePagination"
      :items="paginationItems"
      :total-items="totalItems"
      :items-per-page="query.itemsPerPage"
      :page="query.page"
      class="mt-2 px-4"
      @update:page="updatePage"
      @update:items-per-page="updateItemsPerPage"
      @update:pagination="updatePaginationOptions"
    />
  </div>
</template>

<script>
import { PAGINATION_LIMIT } from '@/config';

import { getPageForNewItemsPerPage } from '@/helpers/pagination';

import AlarmsListTable from './alarms-list-table.vue';

export default {
  components: {
    AlarmsListTable,
  },

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
    stickyHorizontalScroll: {
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
    densable: {
      type: Boolean,
      default: false,
    },
    dense: {
      type: Number,
      default: undefined,
    },
    widgetGridParameters: {
      type: Object,
      default: () => ({}),
    },
    defaultKeepSelectedAfterAction: {
      type: Boolean,
      default: false,
    },
    selectedTags: {
      type: Array,
      default: () => [],
    },
    refreshAlarmsList: {
      type: Function,
      default: () => {},
    },
    eager: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:query', 'update:columns-settings', 'update:dense', 'select:tag', 'remove:tag'],

  computed: {
    totalItems() {
      return this.meta.total_count ?? 0;
    },

    paginationItems() {
      return [5, 10, 20, 50, 100, 200];
    },
  },

  methods: {
    updateQuery(query) {
      this.$emit('update:query', query);
    },

    updatePage(page) {
      this.$emit('update:query', {
        ...this.query,
        page,
      });
    },

    updateItemsPerPage(itemsPerPage) {
      this.$emit('update:query', {
        ...this.query,
        page: getPageForNewItemsPerPage(itemsPerPage, this.query.itemsPerPage, this.query.page),
        itemsPerPage,
      });
    },

    updatePaginationOptions({ page, itemsPerPage }) {
      this.$emit('update:query', {
        ...this.query,
        page,
        itemsPerPage: itemsPerPage || PAGINATION_LIMIT,
      });
    },

    updateColumnsSettings(columnsSettings) {
      this.$emit('update:columns-settings', columnsSettings);
    },
  },
};
</script>
