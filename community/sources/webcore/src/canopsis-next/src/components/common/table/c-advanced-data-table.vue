<template>
  <div class="c-advanced-data-table">
    <v-layout
      wrap
      v-bind="toolbarProps"
    >
      <v-flex
        v-if="shownSearch"
        xs4
      >
        <c-search-field
          v-if="search"
          :value="pagination.search"
          @submit="updateSearchHandler"
          @clear="clearSearchHandler"
        />
        <c-advanced-search-field
          v-else
          :query="pagination"
          :columns="headers"
          :tooltip="searchTooltip"
          @update:query="updatePagination"
        />
      </v-flex>
      <slot
        name="toolbar"
        :selected="selected"
        :update-search="updateSearchHandler"
        :clear-search="clearSearchHandler"
      />
      <v-flex
        v-if="hasMassActionsSlot"
        xs12
      >
        <v-expand-transition>
          <v-layout
            class="px-3 mt-1"
            v-if="selected.length"
          >
            <slot
              name="mass-actions"
              :selected="selected"
              :count="selected.length"
              :clear-selected="clearSelected"
            />
          </v-layout>
        </v-expand-transition>
      </v-flex>
    </v-layout>
    <v-data-table
      v-model="selected"
      :headers="headersWithExpand"
      :items="visibleItems"
      :loading="loading"
      :server-items-length="totalItems"
      :no-data-text="noDataText"
      :options="pagination"
      :header-text="headerText"
      :footer-props="{ itemsPerPageOptions: rowsPerPageItems }"
      :item-key="itemKey"
      :show-select="selectAll"
      :show-expand="expand"
      :is-disabled-item="isDisabledItem"
      :hide-default-footer="hideActions || advancedPagination || noPagination"
      :multi-sort="multiSort"
      :table-class="tableClass"
      :sort-by="disableInitialSort"
      :dense="dense"
      @update:pagination="updatePagination"
    >
      <template #item="props">
        <slot
          v-bind="getItemsProps(props)"
          name="items"
        >
          <tr :key="props.item[itemKey] || props.index">
            <td
              v-if="selectAll || expand"
              @click.stop=""
            >
              <v-layout
                class="c-checkbox-wrapper"
                justify-start
              >
                <slot
                  v-if="selectAll"
                  v-bind="getItemsProps(props)"
                  name="item-select"
                >
                  <v-checkbox-functional
                    v-if="!isDisabledItem(props.item)"
                    v-model="props.selected"
                    primary
                    hide-details
                  />
                  <v-checkbox-functional
                    v-else
                    primary
                    disabled
                    hide-details
                  />
                </slot>
                <slot
                  v-if="expand && isExpandableItem(props.item)"
                  v-bind="getItemsProps(props)"
                  name="item-expand"
                >
                  <c-expand-btn
                    class="ml-2"
                    :expanded="props.expanded"
                    @expand="props.expanded = !props.expanded"
                  />
                </slot>
              </v-layout>
            </td>
            <td
              v-for="header in headers"
              :key="header.value"
            >
              <slot
                :name="header.value"
                v-bind="getItemsProps(props)"
              >
                {{ props.item | get(header.value) }}
              </slot>
            </td>
          </tr>
        </slot>
      </template>
      <template
        v-if="hasExpandSlot"
        #expand="props"
      >
        <div
          class="secondary lighten-2"
          v-if="isExpandableItem(props.item)"
        >
          <slot
            v-bind="props"
            name="expand"
          />
        </div>
      </template>
      <template #headerCell="props">
        <slot
          name="headerCell"
          v-bind="props"
        >
          {{ props.header[headerText] }}
        </slot>
      </template>
      <template #progress="props">
        <slot
          name="progress"
          v-bind="props"
        />
      </template>
    </v-data-table>
    <c-table-pagination
      v-if="!noPagination && pagination && advancedPagination"
      :total-items="totalItems"
      :rows-per-page-items="rowsPerPageItems"
      :rows-per-page="pagination.rowsPerPage"
      :page="pagination.page"
      @update:page="updatePage"
      @update:rows-per-page="updateRecordsPerPage"
    />
  </div>
</template>

<script>
import { omit } from 'lodash';

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
      this.updatePagination({ ...this.pagination, rowsPerPage });
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
