<template lang="pug">
  div.c-advanced-data-table
    v-layout(row, wrap, v-bind="toolbarProps")
      v-flex(v-if="shownSearch", xs4)
        c-search-field(
          v-if="search",
          :value="pagination.search",
          @submit="updateSearchHandler",
          @clear="clearSearchHandler"
        )
        c-advanced-search-field(
          v-else,
          :query="pagination",
          :columns="headers",
          :tooltip="searchTooltip",
          @update:query="updatePagination"
        )
      slot(
        name="toolbar",
        :selected="selected",
        :updateSearch="updateSearchHandler",
        :clearSearch="clearSearchHandler"
      )
      v-flex(v-if="hasMassActionsSlot", xs12)
        v-expand-transition
          v-layout.px-3.mt-1(v-if="selected.length")
            slot(
              name="mass-actions",
              :selected="selected",
              :count="selected.length",
              :clear-selected="clearSelected"
            )
    v-data-table(
      v-model="selected",
      :headers="headersWithExpand",
      :items="visibleItems",
      :loading="loading",
      :total-items="totalItems",
      :no-data-text="noDataText",
      :pagination="pagination",
      :header-text="headerText",
      :rows-per-page-items="rowsPerPageItems",
      :item-key="itemKey",
      :select-all="selectAll",
      :expand="expand",
      :is-disabled-item="isDisabledItem",
      :hide-actions="hideActions || advancedPagination || noPagination",
      :multi-sort="multiSort",
      :table-class="tableClass",
      :disable-initial-sort="disableInitialSort",
      :dense="dense",
      @update:pagination="updatePagination"
    )
      template(#items="props")
        slot(v-bind="getItemsProps(props)", name="items")
          tr(:key="props.item[itemKey] || props.index")
            td(v-if="selectAll || expand", @click.stop="")
              v-layout.c-checkbox-wrapper(row, justify-start)
                slot(v-if="selectAll", v-bind="getItemsProps(props)", name="item-select")
                  v-checkbox-functional(
                    v-if="!isDisabledItem(props.item)",
                    v-model="props.selected",
                    primary,
                    hide-details
                  )
                  v-checkbox-functional(v-else, primary, disabled, hide-details)
                slot(v-if="expand && isExpandableItem(props.item)", v-bind="getItemsProps(props)", name="item-expand")
                  c-expand-btn.ml-2(:expanded="props.expanded", @expand="props.expanded = !props.expanded")
            td(v-for="header in headers", :key="header.value")
              slot(:name="header.value", v-bind="getItemsProps(props)") {{ props.item | get(header.value) }}
      template(v-if="hasExpandSlot", #expand="props")
        div.secondary.lighten-2(v-if="isExpandableItem(props.item)")
          slot(v-bind="props", name="expand")
      template(#headerCell="props")
        slot(name="headerCell", v-bind="props") {{ props.header[headerText] }}
      template(#progress="props")
        slot(name="progress", v-bind="props")
    c-table-pagination(
      v-if="!noPagination && pagination && advancedPagination",
      :total-items="totalItems",
      :rows-per-page-items="rowsPerPageItems",
      :rows-per-page="pagination.rowsPerPage",
      :page="pagination.page",
      @update:page="updatePage",
      @update:rows-per-page="updateRecordsPerPage"
    )
</template>

<script>
import { omit } from 'lodash';

import { getPageForNewRecordsPerPage } from '@/helpers/pagination';

export default {
  model: {
    prop: 'selected',
    event: 'input',
  },
  props: {
    headers: {
      type: Array,
      required: true,
    },
    items: {
      type: Array,
      required: true,
    },
    rowsPerPageItems: {
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
    advancedSearch: {
      type: Boolean,
      default: false,
    },
    advancedPagination: {
      type: Boolean,
      default: false,
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
      required: false,
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
    searchTooltip: {
      type: String,
      default: '',
    },
    pagination: {
      type: Object,
      required: false,
    },
    getPagination: {
      type: Function,
      default: pagination => pagination,
    },
    isDisabledItem: {
      type: Function,
      default: item => !item,
    },
    isExpandableItem: {
      type: Function,
      default: () => true,
    },
    toolbarProps: {
      type: Object,
      required: false,
    },
    multiSort: {
      type: Boolean,
      default: false,
    },
    tableClass: {
      type: String,
      required: false,
    },
    disableInitialSort: {
      type: Boolean,
      default: false,
    },
    dense: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      selectedItems: [],
    };
  },
  computed: {
    selected: {
      get() {
        return this.selectedItems.filter(item => !this.isDisabledItem(item));
      },
      set(selected) {
        this.selectedItems = selected;
      },
    },

    visibleItems() {
      return this.pagination?.rowsPerPage ? this.items.slice(0, this.pagination?.rowsPerPage) : this.items;
    },

    headersWithExpand() {
      if (this.expand && !this.selectAll) {
        return [{ sortable: false, width: 20 }, ...this.headers];
      }

      return this.headers;
    },

    hasExpandSlot() {
      return this.$slots.expand || this.$scopedSlots.expand;
    },

    hasMassActionsSlot() {
      return this.$slots['mass-actions'] || this.$scopedSlots['mass-actions'];
    },

    shownSearch() {
      return this.search || this.advancedSearch;
    },
  },
  watch: {
    items(items) {
      if (this.selectAll) {
        const itemKeys = items.map(item => item[this.itemKey]);

        this.selectedItems = this.selectedItems.filter(selectedItem => itemKeys.includes(selectedItem[this.itemKey]));
      }
    },
  },
  methods: {
    updatePagination(pagination) {
      this.selected = [];

      this.$emit('update:pagination', this.getPagination(pagination));
    },

    updateSearchHandler(search) {
      this.updatePagination({ ...this.pagination, search, page: 1 });
    },

    updateRecordsPerPage(rowsPerPage) {
      this.updatePagination({
        ...this.pagination,

        rowsPerPage,
        page: getPageForNewRecordsPerPage(rowsPerPage, this.pagination.rowsPerPage, this.pagination.page),
      });
    },

    updatePage(page) {
      this.updatePagination({ ...this.pagination, page });
    },

    clearSearchHandler() {
      this.updatePagination(omit(this.pagination, ['search']));
    },

    clearSelected() {
      this.selectedItems = [];
    },

    getItemsProps(state) {
      return {
        item: state.item,
        selected: state.selected,
        disabled: this.isDisabledItem(state.item),
        expanded: state.expanded,
        select: value => state.selected = value || !state.selected,
        expand: value => state.expanded = value || !state.expanded,
      };
    },
  },
};
</script>

<style lang="scss" scoped>
.c-advanced-data-table {
  ::v-deep thead th {
    vertical-align: middle;
  }

  & .c-checkbox-wrapper {
    display: inline-flex;
  }
}
</style>
