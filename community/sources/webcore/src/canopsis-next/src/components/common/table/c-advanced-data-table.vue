<template>
  <div class="c-advanced-data-table">
    <v-layout v-bind="toolbarProps" wrap>
      <v-flex v-if="shownSearch" xs4>
        <c-search
          v-if="search"
          @submit="updateSearchHandler"
        />
        <c-advanced-search
          v-else-if="advancedSearch"
          :columns="headers"
          @submit="updateSearchHandler"
        />
      </v-flex>
      <slot
        :selected="selected"
        :update-search="updateSearchHandler"
        name="toolbar"
      />
      <v-flex
        v-if="hasMassActionsSlot"
        xs12
      >
        <v-expand-transition>
          <v-layout
            v-show="selected.length"
            class="px-2 mt-1"
          >
            <slot
              :selected="selected"
              :count="selected.length"
              :clear-selected="clearSelected"
              name="mass-actions"
            />
          </v-layout>
        </v-expand-transition>
      </v-flex>
    </v-layout>
    <v-data-table
      v-model="selected"
      :headers="preparedHeaders"
      :items="visibleItems"
      :loading="loading"
      :server-items-length="totalItems"
      :no-data-text="noDataText"
      :options="options"
      :header-text="headerText"
      :footer-props="{ itemsPerPageOptions: itemsPerPageItems }"
      :item-key="itemKey"
      :show-select="selectAll"
      :show-expand="expand"
      :item-selectable="isItemSelectable"
      :hide-default-footer="hideActions || advancedPagination || noPagination"
      :table-class="tableClass"
      :dense="dense"
      :loader-height="loaderHeight"
      :ellipsis-headers="ellipsisHeaders"
      checkbox-color="primary"
      @update:options="updateOptions"
    >
      <template v-if="hasItemSlot" #item="props">
        <slot name="item" v-bind="props" />
      </template>

      <template
        v-for="header in headers"
        #[getItemSlotName(header)]="props"
      >
        <slot
          :name="header.value"
          v-bind="getItemsProps(props)"
        >
          {{ props.item | get(header.value) }}
        </slot>
      </template>

      <template
        v-if="hasExpandSlot"
        #expanded-item="props"
      >
        <div
          v-if="isExpandableItem(props.item)"
          class="secondary lighten-2"
        >
          <slot
            v-bind="props"
            name="expand"
          />
        </div>
      </template>

      <template #header="props">
        <slot
          name="header"
          v-bind="props"
        />
      </template>

      <template
        v-for="header in headerScopedSlots"
        #[header]="props"
      >
        <slot
          :name="header"
          v-bind="props"
        />
      </template>

      <template #progress="props">
        <slot
          name="progress"
          v-bind="props"
        />
      </template>
    </v-data-table>

    <c-table-pagination
      v-if="!noPagination && options && advancedPagination"
      :total-items="totalItems"
      :items-per-page="options.itemsPerPage"
      :items="itemsPerPageItems"
      :page="options.page"
      @update:page="updatePage"
      @update:items-per-page="updateItemsPerPage"
    />
  </div>
</template>

<script>
import { getPageForNewItemsPerPage } from '@/helpers/pagination';

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
      default: () => true,
    },
    toolbarProps: {
      type: Object,
      required: false,
    },
    tableClass: {
      type: String,
      required: false,
    },
    dense: {
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
  },
  data() {
    return {
      selectedItems: [],
    };
  },
  computed: {
    preparedHeaders() {
      const headers = [];

      if (this.selectAll) {
        headers.push({
          value: 'data-table-select',
          text: '',
          sortable: false,
          width: '1px',
        });
      }

      if (this.expand) {
        headers.push({
          value: 'data-table-expand',
          text: '',
          sortable: false,
          width: '1px',
        });
      }

      return [
        ...headers,
        ...this.headers,
      ];
    },

    headerScopedSlots() {
      return Object.keys(this.$scopedSlots ?? {}).filter(name => name.startsWith('header.'));
    },

    selected: {
      get() {
        return this.selectedItems.filter(item => !this.isDisabledItem(item));
      },
      set(selected) {
        this.selectedItems = selected;
      },
    },

    visibleItems() {
      return this.options?.itemsPerPage ? this.items.slice(0, this.options?.itemsPerPage) : this.items;
    },

    hasItemSlot() {
      return this.$slots.item || this.$scopedSlots.item;
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
    updateOptions(options) {
      this.selected = [];

      this.$emit('update:options', options);
    },

    updateSearchHandler(search) {
      this.updateOptions({ ...this.options, search, page: 1 });
    },

    updateItemsPerPage(itemsPerPage) {
      this.updateOptions({
        ...this.options,

        itemsPerPage,
        page: getPageForNewItemsPerPage(itemsPerPage, this.options.itemsPerPage, this.options.page),
      });
    },

    updatePage(page) {
      this.updateOptions({ ...this.options, page });
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
        expand: state.expand,
      };
    },

    getItemSlotName(header) {
      return `item.${header.value}`;
    },

    isItemSelectable(item) {
      return !this.isDisabledItem(item);
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
