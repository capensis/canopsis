<template>
  <div class="c-advanced-data-table">
    <v-row v-bind="toolbarProps" class="px-4" wrap>
      <div v-if="shownSearch" class="c-advanced-data-table__search">
        <c-search
          v-if="search"
          :model-value="options?.search"
          :label="searchLabel"
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
          v-if="checkIsExpandableItem(props.item)"
          :class="{ 'v-data-table__expand-icon--active': isItemExpanded(props.item) }"
          class="v-data-table__expand-icon"
          @click.stop="toggleExpand(props.item)"
        >
          {{ expandIcon }}
        </c-icon>
        <span v-else />
      </template>

      <template v-if="hasExpandSlot" #expanded-row="props">
        <tr>
          <td
            :colspan="preparedHeaders.length"
            :class="{
              'c-advanced-data-table__expand-cell--closing': isItemClosing(props.item),
            }"
            class="c-advanced-data-table__expand-cell"
          >
            <div
              :ref="(el) => setExpandElementRef(el, props.item)"
              :class="{
                'c-advanced-data-table__expand--closing': isItemClosing(props.item),
              }"
              class="c-advanced-data-table__expand"
            >
              <div :class="expandContentClass">
                <slot v-bind="getItemsProps(props)" name="expand" />
              </div>
            </div>
          </td>
        </tr>
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
      @update:page="updatePage"
      @update:items-per-page="updateItemsPerPage"
      @update:pagination="updatePaginationOptions"
    />
  </div>
</template>

<script>
import get from 'lodash/get'
import isEqual from 'lodash/isEqual'
import pick from 'lodash/pick'

import { mapIds } from '@/helpers/array'
import { convertDataTableSortToQuery } from '@/helpers/entities/shared/query'
import {
  prepareQueryWithAdvancedSearch,
  prepareQueryWithoutAdvancedSearch,
} from '@/helpers/search/advanced-search'

const EXPAND_TRANSITION_DURATION = 260

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
      default: (item) => !item,
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
      default: 'secondary lighten-2 pa-4',
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
      openingExpandedKeys: [],
      closingExpandedKeys: [],
      expandElements: {},
      expandTimeouts: {},
      keepSelectedAfterAction: this.defaultKeepSelectedAfterAction ?? false,
    }
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
      ])
    },

    vuetifySortBy() {
      const sortBy = this.vDataTableOptions.sortBy ?? []
      const sortDesc = this.vDataTableOptions.sortDesc ?? []

      return sortBy.map((item, index) =>
        item.key
          ? item
          : {
              key: item,
              order: sortDesc[index] ? 'desc' : 'asc',
            },
      )
    },

    preparedHeaders() {
      const headers = []

      if (this.selectAll) {
        headers.push({
          key: 'data-table-select',
          title: '',
          sortable: false,
          width: '1px',
        })
      }

      if (this.expand) {
        headers.push({
          key: 'data-table-expand',
          title: '',
          sortable: false,
          width: '1px',
        })
      }

      return [
        ...headers,
        ...this.headers.map((header) => ({
          ...header,
          key: header.key ?? header.value,
          title: header.title ?? header[this.headerText],
        })),
      ]
    },

    headerScopedSlots() {
      return Object.keys(this.$slots ?? {}).filter(
        (name) => name.startsWith('header.') && name !== 'header.data-table-select',
      )
    },

    selected: {
      get() {
        return this.selectedItems.filter((item) => !this.isDisabledItem(item))
      },
      set(selected) {
        this.selectedItems = selected
      },
    },

    selectedKeys() {
      return mapIds(this.selected, this.itemKey)
    },

    hasHeaderSelectSlot() {
      return this.$slots['header.data-table-select']
    },

    hasItemSlot() {
      return this.$slots.item
    },

    hasExpandSlot() {
      return this.$slots.expand
    },

    hasHeaderSlot() {
      return this.$slots.header
    },

    hasProgressSlot() {
      return this.$slots.progress
    },

    hasBodyPrependSlot() {
      return this.$slots['body.prepend']
    },

    hasMassActionsSlot() {
      return this.$slots['mass-actions']
    },

    shownSearch() {
      return this.search || this.advancedSearch
    },

    density() {
      return {
        [this.ultraDense]: 'compact',
        [this.dense]: 'comfortable',
      }.true
    },

    itemsLength() {
      return this.totalItems || this.items.length
    },
  },

  watch: {
    items(items) {
      const itemKeys = new Set(items.map((item) => this.getItemKey(item)))

      if (this.selectAll) {
        const selectedKeys = new Set(this.selectedItems.map((item) => this.getItemKey(item)))

        this.selectedItems = items.filter(
          (item) => selectedKeys.has(this.getItemKey(item)) && !this.isDisabledItem(item),
        )
      }

      const expandedKeys = new Set(this.expandedItems.map((item) => this.getItemKey(item)))

      this.expandedItems = items.filter(
        (item) => expandedKeys.has(this.getItemKey(item)) && this.checkIsExpandableItem(item),
      )
      this.openingExpandedKeys = this.openingExpandedKeys.filter((key) => itemKeys.has(key))
      this.closingExpandedKeys = this.closingExpandedKeys.filter((key) => itemKeys.has(key))
    },
  },

  beforeUnmount() {
    Object.values(this.expandTimeouts).forEach(clearTimeout)
  },

  methods: {
    getItemKey(item) {
      return item?.[this.itemKey]
    },

    checkIsExpandableItem(item) {
      return this.isExpandableItem ? this.isExpandableItem(item) : true
    },

    setExpandElementRef(element, item) {
      const key = this.getItemKey(item)

      if (!key) {
        return
      }

      if (!element) {
        delete this.expandElements[key]

        return
      }

      this.expandElements[key] = element

      if (this.openingExpandedKeys.includes(key)) {
        this.$nextTick(() => this.openExpandElement(key))
      }
    },

    openExpandElement(key) {
      const element = this.expandElements[key]

      if (!element) {
        return
      }

      element.style.height = '0px'

      requestAnimationFrame(() => {
        element.style.height = `${element.scrollHeight}px`
      })
    },

    closeExpandElement(key) {
      const element = this.expandElements[key]

      if (!element) {
        return
      }

      element.style.height = `${element.scrollHeight}px`

      requestAnimationFrame(() => {
        element.style.height = '0px'
      })
    },

    isItemExpanded(item) {
      const key = this.getItemKey(item)

      return (
        this.expandedItems.some((expandedItem) => this.getItemKey(expandedItem) === key) &&
        !this.closingExpandedKeys.includes(key)
      )
    },

    isItemOpening(item) {
      return this.openingExpandedKeys.includes(this.getItemKey(item))
    },

    isItemClosing(item) {
      return this.closingExpandedKeys.includes(this.getItemKey(item))
    },

    toggleExpand(item) {
      if (!this.checkIsExpandableItem(item)) {
        return
      }

      const key = this.getItemKey(item)

      if (this.expandTimeouts[key]) {
        clearTimeout(this.expandTimeouts[key])
        delete this.expandTimeouts[key]
      }

      if (this.isItemExpanded(item)) {
        this.openingExpandedKeys = this.openingExpandedKeys.filter(
          (openingKey) => openingKey !== key,
        )
        this.closingExpandedKeys = [...this.closingExpandedKeys, key]

        this.$nextTick(() => this.closeExpandElement(key))

        this.expandTimeouts[key] = window.setTimeout(() => {
          this.expandedItems = this.expandedItems.filter(
            (expandedItem) => this.getItemKey(expandedItem) !== key,
          )
          this.closingExpandedKeys = this.closingExpandedKeys.filter(
            (closingKey) => closingKey !== key,
          )

          delete this.expandTimeouts[key]
        }, EXPAND_TRANSITION_DURATION)

        return
      }

      this.closingExpandedKeys = this.closingExpandedKeys.filter((closingKey) => closingKey !== key)
      this.openingExpandedKeys = [...this.openingExpandedKeys, key]

      if (!this.expandedItems.some((expandedItem) => this.getItemKey(expandedItem) === key)) {
        this.expandedItems = [...this.expandedItems, item]
      }

      this.expandTimeouts[key] = window.setTimeout(() => {
        const element = this.expandElements[key]

        if (element) {
          element.style.height = 'auto'
        }

        this.openingExpandedKeys = this.openingExpandedKeys.filter(
          (openingKey) => openingKey !== key,
        )

        delete this.expandTimeouts[key]
      }, EXPAND_TRANSITION_DURATION)
    },

    updateOptions(options) {
      if (isEqual(this.options, options)) {
        return
      }

      this.selected = []
      this.expandedItems = []
      this.openingExpandedKeys = []
      this.closingExpandedKeys = []

      Object.values(this.expandTimeouts).forEach(clearTimeout)
      this.expandTimeouts = {}

      this.$emit('update:options', options)
    },

    updateDataTableOptions(options) {
      const sort = convertDataTableSortToQuery(options.sortBy)

      return this.updateOptions({
        ...this.options,
        ...pick(options, ['page', 'itemsPerPage', 'groupBy', 'mustSort', 'multiSort']),
        sortBy: sort.sortBy,
        sortDesc: sort.sortDesc,
      })
    },

    updateSearch(search) {
      this.updateOptions({ ...this.options, search, page: 1 })
    },

    updateAdvancedSearch(search) {
      this.updateOptions(prepareQueryWithAdvancedSearch(this.options, search))
    },

    resetSearch() {
      this.updateOptions(prepareQueryWithoutAdvancedSearch(this.options))
    },

    updateItemsPerPage(itemsPerPage) {
      this.updateOptions({ ...this.options, itemsPerPage })
    },

    updatePage(page) {
      this.updateOptions({ ...this.options, page })
    },

    updatePaginationOptions({ page, itemsPerPage }) {
      this.updateOptions({ ...this.options, page, itemsPerPage })
    },

    clearSelected(keepSelected = this.keepSelectedAfterAction, force = false) {
      if (keepSelected && !force) {
        return
      }

      this.selectedItems = []
    },

    getItemsProps(state) {
      return {
        index: state.index,
        item: state.item,
        selected: state.selected,
        disabled: this.isDisabledItem(state.item),
        expanded: this.isItemExpanded(state.item),
        select: state.toggleSelect,
        expand: this.toggleExpand,
      }
    },

    getHeaderKey(header) {
      return header.value ?? header.key
    },

    getItemSlotName(header) {
      return `item.${this.getHeaderKey(header)}`
    },

    getItemValue(item, value) {
      return get(item, value)
    },

    isItemSelectable(item) {
      return !this.isDisabledItem(item)
    },
  },
}
</script>

<style lang="scss" scoped>
.c-advanced-data-table {
  :deep(thead th) {
    vertical-align: middle;
  }

  &__expand-cell {
    height: 0;
    padding: 0 !important;
    border-bottom: 0 !important;
    line-height: 0;
  }

  &__search {
    display: flex;
    flex: 0 0 33.333333%;
    max-width: 33.333333%;
  }

  &__expand {
    height: auto;
    overflow: hidden;
    opacity: 1;
    line-height: normal;
    transition:
      height 0.2s ease-out,
      opacity 0.2s ease-out;
  }

  &__expand--opening,
  &__expand--closing {
    will-change: height, opacity;
  }

  &__expand--opening {
    opacity: 1;
  }

  &__expand--closing {
    opacity: 0;
  }

  &__expand-cell--closing {
    padding: 0 !important;
    border-bottom: 0 !important;
  }

  & .c-checkbox-wrapper {
    display: inline-flex;
  }
}
</style>
