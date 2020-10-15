<template lang="pug">
  div.white
    v-layout(row, wrap)
      v-flex(v-if="shownSearch", xs4)
        search-field(
          v-if="search",
          :value="pagination.search",
          @submit="updateSearchHandler",
          @clear="clearSearchHandler"
        )
        advanced-search(
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
    v-data-table(
      v-model="selected",
      :headers="headers",
      :items="items",
      :loading="loading",
      :total-items="totalItems",
      :pagination="pagination",
      :header-text="headerText",
      :item-key="itemKey",
      :select-all="selectAll",
      :expand="expand",
      :is-disabled-item="isDisabledItem",
      :hide-actions="hideActions || advancedPagination",
      @update:pagination="updatePagination($event)"
    )
      template(slot="items", slot-scope="props")
        slot(v-bind="getItemsProps(props)", name="items")
          tr(:key="props.item[itemKey] || props.index")
            td(v-if="selectAll || expand", @click.stop)
              v-layout.checkbox-wrapper(row, justify-start)
                slot(v-if="selectAll", v-bind="getItemsProps(props)", name="item-select")
                  v-checkbox-functional(
                    v-if="!isDisabledItem(props.item)",
                    v-model="props.selected",
                    primary,
                    hide-details
                  )
                  v-checkbox-functional(v-else, primary, disabled, hide-details)
                slot(v-if="expand", name="item-expand", v-bind="getItemsProps(props)")
                  expand-button.ml-2(:expanded="props.expanded", @expand="props.expanded = !props.expanded")
            td(v-for="header in headers", :key="header.value")
              slot(:name="header.value", v-bind="getItemsProps(props)") {{ props.item | get(header.value) }}
      template(v-if="hasExpandSlot", slot="expand", slot-scope="props")
        div.secondary.lighten-2
          slot(v-bind="props", name="expand")
      template(slot="headerCell", slot-scope="props")
        slot(name="headerCell", v-bind="props") {{ props.header[headerText] }}
    slot(name="mass-actions", :selected="selected")
    v-layout.white(v-show="totalItems && advancedPagination", align-center)
      v-flex(xs10)
        pagination(
          :page="pagination.page",
          :limit="pagination.rowsPerPage",
          :total="totalItems",
          @input="updatePage"
        )
      v-spacer
      v-flex(xs2)
        records-per-page(:value="pagination.rowsPerPage", @input="updateRecordsPerPage")
</template>

<script>
import { omit } from 'lodash';

import SearchField from '@/components/forms/fields/search-field.vue';
import AdvancedSearch from '@/components/other/shared/search/advanced-search.vue';
import ExpandButton from '@/components/other/buttons/expand-button.vue';
import Pagination from '@/components/tables/pagination.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';

export default {
  components: {
    SearchField,
    AdvancedSearch,
    ExpandButton,
    Pagination,
    RecordsPerPage,
  },
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
    totalItems: {
      type: Number,
      required: false,
    },
    itemKey: {
      type: String,
      default: '_id',
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
      required: true,
    },
    getPagination: {
      type: Function,
      default: pagination => pagination,
    },
    isDisabledItem: {
      type: Function,
      default: item => !item,
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
    hasExpandSlot() {
      return this.$slots.expand || this.$scopedSlots.expand;
    },

    shownSearch() {
      return this.search || this.advancedSearch;
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
  .checkbox-wrapper {
    display: inline-flex;
  }
</style>
