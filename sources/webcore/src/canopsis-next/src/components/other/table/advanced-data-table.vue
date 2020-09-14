<template lang="pug">
  div.white
    v-layout(row, wrap)
      v-flex(v-if="shownSearch", xs4)
        search-field(v-if="search", @submit="updateSearchHandler", @clear="clearSearchHandler")
        advanced-search(
          v-else,
          :query="pagination",
          :columns="headers",
          :tooltip="searchTooltip",
          @update:query="updatePagination"
        )
      slot(name="toolbar")
    v-data-table(
      v-field="selected",
      :headers="headers",
      :items="items",
      :loading="loading",
      :total-items="totalItems",
      :pagination="pagination",
      :header-text="headerText",
      :item-key="itemKey",
      :select-all="selectAll",
      :expand="expand",
      :hide-actions="hideActions",
      @update:pagination="updatePagination($event)"
    )
      template(slot="items", slot-scope="props")
        slot(v-bind="getItemsProps(props)", name="items")
          tr(:key="props.item[itemKey] || props.index", @click="expandPanel(props)")
            td(v-if="selectAll", @click.stop)
              slot(name="selectAll", v-bind="getItemsProps(props)")
                v-checkbox(v-model="props.selected", hide-details)
            td(v-for="header in headers", :key="header.value")
              slot(:name="header.value", v-bind="getItemsProps(props)") {{ props.item | get(header.value) }}
      template(v-if="hasExpandSlot", slot="expand", slot-scope="props")
        div.secondary.lighten-2
          slot(name="expand", v-bind="props")
      template(slot="headerCell", slot-scope="props")
        slot(name="headerCell", v-bind="props") {{ props.header[headerText] }}
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
import Pagination from '@/components/tables/pagination.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';

export default {
  components: {
    SearchField,
    AdvancedSearch,
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
    selected: {
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
  },
  computed: {
    hasExpandSlot() {
      return this.$slots.expand || this.$scopedSlots.expand;
    },

    shownSearch() {
      return this.search || this.advancedSearch;
    },
  },
  methods: {
    updatePagination(pagination) {
      this.$emit('update:pagination', this.getPagination(pagination));
    },

    updateSearchHandler(search) {
      this.updatePagination({ ...this.pagination, search });
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
        expanded: state.expanded,
        select: value => state.selected = value || !state.selected,
        expand: value => state.expanded = value || !state.expanded,
      };
    },

    expandPanel(state) {
      if (this.expand) {
        state.expanded = !state.expanded;
      }
    },
  },
};
</script>
