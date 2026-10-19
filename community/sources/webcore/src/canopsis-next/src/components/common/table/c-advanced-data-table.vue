<template>
  <div class="c-advanced-data-table">
    <v-row v-bind="toolbarProps" class="px-4" wrap>
      <div v-if="shownSearch" class="c-advanced-data-table__search">
        <c-search
          v-if="search"
          :model-value="options?.search"
          :label="searchLabel"
          :disabled="searchDisabled"
          @submit="updateSearch"
        />
        <c-advanced-search
          v-else-if="advancedSearch"
          :attributes="advancedSearchAttributes"
          :pending="advancedSearchLoading"
          @submit="updateAdvancedSearch"
          @reset="resetSearch"
        />
      </div>

      <slot
        :selected="selected"
        :update-search="updateSearch"
        :update-advanced-search="updateAdvancedSearch"
        :reset-search="resetSearch"
        name="toolbar"
      />

      <c-mass-actions-panel
        v-if="selectAll && selected.length && hasMassActionsSlot && !hideMassActions"
        v-model="keepSelectedAfterAction"
        v-bind="massActionsGridParameters"
        :selected="selected"
        :default-value="defaultKeepSelectedAfterAction"
        @clear:selected="clearSelected(keepSelectedAfterAction, true)"
      >
        <template #actions>
          <slot
            :selected="selected"
            :selected-keys="selectedKeys"
            :count="selected.length"
            :clear-selected="clearSelected"
            name="mass-actions"
          />
        </template>
      </c-mass-actions-panel>
    </v-row>

    <v-data-table-server
      v-model="selected"
      v-model:expanded="expandedItems"
      ref="dataTable"
      :headers="preparedHeaders"
      :items="items"
      :loading="loading"
      :items-length="itemsLength"
      :no-data-text="noDataText"
      :page="vDataTableOptions.page"
      :items-per-page="vDataTableOptions.itemsPerPage"
      :sort-by="vuetifySortBy"
      :item-value="itemKey"
      :show-select="selectAll"
      :show-expand="expand"
      :item-selectable="isItemSelectable"
      :hide-default-footer="hideActions || advancedPagination || noPagination"
      :class="tableClass"
      :density="density"
      :loader-height="loaderHeight"
      :expand-icon="expandIcon"
      color="primary"
      return-object
      @update:options="updateDataTableOptions"
    >
      <template v-if="hasHeaderSelectSlot" #[`header.data-table-select`]="props">
        <slot name="header.data-table-select" v-bind="props" />
      </template>

      <template v-if="hasItemSlot" #item="props">
        <slot name="item" v-bind="getItemsProps(props)" />
      </template>

      <template v-for="header in headers" #[getItemSlotName(header)]="props">
        <slot :name="getHeaderKey(header)" v-bind="getItemsProps(props)">
          {{ getItemValue(props.item, getHeaderKey(header)) }}
        </slot>
      </template>

      <template #[`item.data-table-expand`]="props">
        <c-icon
          v-if="checkIsExpandableItem(getRawItem(props.item))"
          :class="{ 'v-data-table__expand-icon--active': isItemExpanded(getRawItem(props.item)) }"
          class="v-data-table__expand-icon"
          @click.stop="toggleExpand(getRawItem(props.item))"
        >
          {{ expandIcon }}
        </c-icon>
        <span v-else />
      </template>

      <template v-if="hasExpandSlot" #expanded-row="props">
        <c-table-expand-row
          :ref="(el) => setExpandRowRef(el, props.item)"
          :colspan="preparedHeaders.length"
          :content-class="expandContentClass"
          expanded
          @after-close="removeExpandedItem(props.item)"
        >
          <slot v-bind="getItemsProps(props)" name="expand" />
        </c-table-expand-row>
      </template>

      <template v-if="hasHeaderSlot" #headers="props">
        <slot name="header" v-bind="props" />
      </template>

      <template v-for="header in headerScopedSlots" #[header]="props">
        <slot :name="header" v-bind="props" />
      </template>

      <template v-if="hasProgressSlot" #loader="props">
        <slot name="progress" v-bind="props" />
      </template>

      <template v-if="hasBodyPrependSlot" #[`body.prepend`]="props">
        <slot name="body.prepend" v-bind="props" />
      </template>
    </v-data-table-server>

    <c-table-pagination
      v-if="!noPagination && options && advancedPagination"
      :total-items="itemsLength"
      :items-per-page="options.itemsPerPage"
      :items="itemsPerPageItems"
      :page="options.page"
      :total-visible="paginationTotalVisible"
      @update:page="updatePage"
      @update:items-per-page="updateItemsPerPage"
      @update:pagination="updatePaginationOptions"
    />
  </div>
</template>

<script>
import get from 'lodash/get';
import isEqual from 'lodash/isEqual';
import pick from 'lodash/pick';

import { mapIds } from '@/helpers/array';
import { convertDataTableSortToQuery } from '@/helpers/entities/shared/query';
import {
  prepareQueryWithAdvancedSearch,
  prepareQueryWithoutAdvancedSearch,
} from '@/helpers/search/advanced-search';

export default {
  props: {
    headers: {
      type: Array,
      required: true,
    },
    items: {
      type: Array,
      required: true,
    },
    itemsPerPageItems: {
      type: Array,
      required: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    selectAll: {
      type: Boolean,
      default: false,
    },
    expand: {
      type: Boolean,
      default: false,
    },
    search: {
      type: Boolean,
      default: false,
    },
    searchDisabled: {
      type: Boolean,
      default: false,
    },
    advancedSearch: {
      type: Boolean,
      default: false,
    },
    advancedPagination: {
      type: Boolean,
      default: false,
    },
    paginationTotalVisible: {
      type: Number,
      required: false,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
    noPagination: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      default: 0,
    },
    itemKey: {
      type: String,
      default: '_id',
    },
    noDataText: {
      type: String,
      required: false,
    },
    headerText: {
      type: String,
      default: 'text',
    },
    options: {
      type: Object,
      required: false,
    },
    isDisabledItem: {
      type: Function,
      default: item => !item,
    },
    isExpandableItem: {
      type: Function,
      required: false,
    },
    toolbarProps: {
      type: Object,
      required: false,
    },
    searchLabel: {
      type: String,
      default: '',
    },
    tableClass: {
      type: [String, Object],
      required: false,
    },
    dense: {
      type: Boolean,
      default: false,
    },
    ultraDense: {
      type: Boolean,
      default: false,
    },
    ellipsisHeaders: {
      type: Boolean,
      default: false,
    },
    loaderHeight: {
      type: [String, Number],
      default: 2,
    },
    expandIcon: {
      type: String,
      default: '$expand',
    },
    expandContentClass: {
      type: [String, Array, Object],
      default: 'secondary lighten-2',
    },
    advancedSearchAttributes: {
      type: Array,
      default: () => [],
    },
    advancedSearchLoading: {
      type: Boolean,
      default: false,
    },
    massActionsGridParameters: {
      type: Object,
      required: false,
    },
    defaultKeepSelectedAfterAction: {
      type: Boolean,
      default: false,
    },
    hideMassActions: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:options'],

  data() {
    return {
      selectedItems: [],
      expandedItems: [],
      expandRows: {},
      keepSelectedAfterAction: this.defaultKeepSelectedAfterAction ?? false,
    };
  },

  computed: {
    vDataTableOptions() {
      return pick(this.options, [
        'page',
        'itemsPerPage',
        'sortBy',
        'sortDesc',
        'groupBy',
        'groupDesc',
        'mustSort',
        'multiSort',
      ]);
    },

    vuetifySortBy() {
      const sortBy = this.vDataTableOptions.sortBy ?? [];
      const sortDesc = this.vDataTableOptions.sortDesc ?? [];

      return sortBy.map((item, index) => item.key
          ? item
          : {
              key: item,
              order: sortDesc[index] ? 'desc' : 'asc',
            },
      );
    },

    preparedHeaders() {
      const headers = [];

      if (this.selectAll) {
        headers.push({
          key: 'data-table-select',
          title: '',
          sortable: false,
          width: '1px',
        });
      }

      if (this.expand) {
        headers.push({
          key: 'data-table-expand',
          title: '',
          sortable: false,
          width: '1px',
          headerProps: {
            noPadding: false,
          },
          cellProps: {
            noPadding: false,
          },
        });
      }

      return [
        ...headers,
        ...this.headers.map(header => ({
          ...header,
          key: header.key ?? header.value,
          title: header.title ?? header[this.headerText],
        })),
      ];
    },

    headerScopedSlots() {
      return Object.keys(this.$slots ?? {}).filter(
        name => name.startsWith('header.') && name !== 'header.data-table-select',
      );
    },

    selected: {
      get() {
        return this.selectedItems.filter(item => !this.isDisabledItem(item));
      },
      set(selected) {
        this.selectedItems = selected;
      },
    },

    selectedKeys() {
      return mapIds(this.selected, this.itemKey);
    },

    hasHeaderSelectSlot() {
      return this.$slots['header.data-table-select'];
    },

    hasItemSlot() {
      return this.$slots.item;
    },

    hasExpandSlot() {
      return this.$slots.expand;
    },

    hasHeaderSlot() {
      return this.$slots.header;
    },

    hasProgressSlot() {
      return this.$slots.progress;
    },

    hasBodyPrependSlot() {
      return this.$slots['body.prepend'];
    },

    hasMassActionsSlot() {
      return this.$slots['mass-actions'];
    },

    shownSearch() {
      return this.search || this.advancedSearch;
    },

    density() {
      return {
        [this.ultraDense]: 'compact',
        [this.dense]: 'comfortable',
      }.true;
    },

    itemsLength() {
      return this.totalItems || this.items.length;
    },
  },

  watch: {
    items(items) {
      if (this.selectAll) {
        const selectedKeys = new Set(this.selectedItems.map(item => this.getItemKey(item)));

        this.selectedItems = items.filter(
          item => selectedKeys.has(this.getItemKey(item)) && !this.isDisabledItem(item),
        );
      }

      const expandedKeys = new Set(this.expandedItems.map(item => this.getItemKey(item)));

      this.expandedItems = items.filter(
        item => expandedKeys.has(this.getItemKey(item)) && this.checkIsExpandableItem(item),
      );
    },
  },

  methods: {
    getRawItem(item) {
      return item?.raw ?? item;
    },

    getItemKey(item) {
      return get(this.getRawItem(item), this.itemKey);
    },

    checkIsExpandableItem(item) {
      return this.isExpandableItem ? this.isExpandableItem(item) : true;
    },

    setExpandRowRef(element, item) {
      const key = this.getItemKey(item);

      if (!key) {
        return;
      }

      if (!element) {
        delete this.expandRows[key];

        return;
      }

      this.expandRows[key] = element;
    },

    isItemExpanded(item) {
      const key = this.getItemKey(item);

      return this.expandedItems.some(expandedItem => this.getItemKey(expandedItem) === key);
    },

    toggleExpand(item) {
      if (!this.checkIsExpandableItem(item)) {
        return;
      }

      const key = this.getItemKey(item);

      if (this.isItemExpanded(item)) {
        this.expandRows[key]?.hide();

        return;
      }

      if (!this.expandedItems.some(expandedItem => this.getItemKey(expandedItem) === key)) {
        this.expandedItems = [...this.expandedItems, item];
      }
    },

    removeExpandedItem(item) {
      const key = this.getItemKey(item);

      this.expandedItems = this.expandedItems.filter(
        expandedItem => this.getItemKey(expandedItem) !== key,
      );
    },

    updateOptions(options) {
      if (isEqual(this.options, options)) {
        return;
      }

      this.selected = [];
      this.expandedItems = [];
      this.expandRows = {};

      this.$emit('update:options', options);
    },

    updateDataTableOptions(options) {
      const sort = convertDataTableSortToQuery(options.sortBy);

      return this.updateOptions({
        ...this.options,
        ...pick(options, ['page', 'itemsPerPage', 'groupBy', 'mustSort', 'multiSort']),
        sortBy: sort.sortBy,
        sortDesc: sort.sortDesc,
      });
    },

    updateSearch(search) {
      this.updateOptions({ ...this.options, search, page: 1 });
    },

    updateAdvancedSearch(search) {
      this.updateOptions(prepareQueryWithAdvancedSearch(this.options, search));
    },

    resetSearch() {
      this.updateOptions(prepareQueryWithoutAdvancedSearch(this.options));
    },

    updateItemsPerPage(itemsPerPage) {
      this.updateOptions({ ...this.options, itemsPerPage });
    },

    updatePage(page) {
      this.updateOptions({ ...this.options, page });
    },

    updatePaginationOptions({ page, itemsPerPage }) {
      this.updateOptions({ ...this.options, page, itemsPerPage });
    },

    clearSelected(keepSelected = this.keepSelectedAfterAction, force = false) {
      if (keepSelected && !force) {
        return;
      }

      this.selectedItems = [];
    },

    getItemsProps(state) {
      const item = this.getRawItem(state.item);

      return {
        index: state.index,
        item,
        selected: state.selected,
        disabled: this.isDisabledItem(item),
        expanded: this.isItemExpanded(item),
        select: state.toggleSelect,
        expand: this.toggleExpand,
      };
    },

    getHeaderKey(header) {
      return header.value ?? header.key;
    },

    getItemSlotName(header) {
      return `item.${this.getHeaderKey(header)}`;
    },

    getItemValue(item, value) {
      return get(item, value);
    },

    isItemSelectable(item) {
      return !this.isDisabledItem(item);
    },
  },
};
</script>

<style lang="scss" scoped>
.c-advanced-data-table {
  :deep(thead th) {
    vertical-align: middle;
  }

  :deep(.v-data-table__expand-icon) {
    cursor: pointer;
    user-select: none;
    transition: transform 0.2s ease-out;
  }

  :deep(.v-data-table__expand-icon--active) {
    transform: rotate(-180deg);
  }

  &__search {
    display: flex;
    align-items: center;
    flex: 0 0 33.333333%;
    max-width: 33.333333%;

    :deep(.c-action-btn) {
      align-self: center;
    }
  }

  & .c-checkbox-wrapper {
    display: inline-flex;
  }
}
</style>

<style lang="scss">
.c-advanced-data-table__search {
  flex: 1 1 0;
  max-width: none;
}

.c-advanced-data-table__toolbar-field {
  display: flex;
  align-items: center;
  flex: 0 0 auto;

  .c-field {
    width: max-content;
    flex: 0 0 auto;
  }
}
</style>
