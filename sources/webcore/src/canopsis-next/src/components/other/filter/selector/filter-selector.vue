<template lang="pug">
  v-layout(align-center, row, wrap)
    v-flex(v-show="!hideSelect", v-bind="flexProps.switch")
      v-switch(
      :label="$t('filterSelector.fields.mixFilters')",
      :input-value="isMultiple",
      :disabled="!hasAccessToEditFilter",
      @change="updateIsMultipleFlag"
      )
    v-flex(v-show="!hideSelect && isMultiple", v-bind="flexProps.radio")
      v-radio-group(
      :value="condition",
      :disabled="!hasAccessToEditFilter",
      @change="updateCondition"
      )
        v-radio(label="AND", value="$and")
        v-radio(label="OR", value="$or")
    v-flex(v-show="!hideSelect", v-bind="flexProps.select")
      v-select(
      :value="value",
      :items="preparedFilters",
      :label="label",
      :itemText="itemText",
      :itemValue="itemValue",
      :multiple="isMultiple",
      :disabled="!hasAccessToEditFilter",
      return-object,
      clearable,
      @input="updateSelectedFilter"
      )
        template(slot="item", slot-scope="{ parent, item, tile }")
          v-list-tile-action(v-if="isMultiple", @click.stop="parent.$emit('select', item)")
            v-checkbox(:inputValue="tile.props.value", :color="parent.color")
          v-list-tile-content
            v-list-tile-title
              span {{ item[itemText] }}
              v-icon.ml-2(
              v-show="!hideSelectIcon",
              :color="tile.props.value ? parent.color : ''",
              small
              ) {{ item.locked ? 'lock' : 'person' }}
    v-flex(v-bind="flexProps.list")
      v-btn(v-if="!long", @click="showFiltersListModal", icon, small)
        v-icon filter_list
      filters-list(
      v-else,
      :filters="filters",
      :hasAccessToAddFilter="hasAccessToAddFilter",
      :hasAccessToEditFilter="hasAccessToEditFilter",
      @create:filter="createFilter",
      @update:filter="updateFilter",
      @delete:filter="deleteFilter"
      )
</template>

<script>
import { MODALS, FILTER_DEFAULT_VALUES } from '@/constants';

import modalMixin from '@/mixins/modal';

import FiltersList from '@/components/other/filter/list/filters-list.vue';

export default {
  components: { FiltersList },
  mixins: [modalMixin],
  props: {
    long: {
      type: Boolean,
      default: false,
    },
    value: {
      type: [Object, Array],
      default: () => null,
    },
    filters: {
      type: Array,
      default: () => [],
    },
    lockedFilters: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    itemText: {
      type: String,
      default: 'title',
    },
    itemValue: {
      type: String,
      default: 'filter',
    },
    condition: {
      type: String,
      default: FILTER_DEFAULT_VALUES.condition,
    },
    hideSelect: {
      type: Boolean,
      default: false,
    },
    hideSelectIcon: {
      type: Boolean,
      default: false,
    },
    hasAccessToAddFilter: {
      type: Boolean,
      default: true,
    },
    hasAccessToEditFilter: {
      type: Boolean,
      default: true,
    },
  },
  computed: {
    flexProps() {
      return {
        switch: this.long ? { xs6: true } : {},
        radio: this.long ? { xs6: true } : {},
        select: this.long ? { xs12: true } : {},
        list: this.long ? { xs12: true } : {},
      };
    },

    isMultiple() {
      return Array.isArray(this.value);
    },

    preparedFilters() {
      const preparedFilters = [...this.filters];

      if (this.lockedFilters.length) {
        return preparedFilters.concat(
          { divider: true },
          this.lockedFilters.map(filter => ({ ...filter, locked: true })),
        );
      }

      return preparedFilters;
    },
  },
  methods: {
    updateIsMultipleFlag(value) {
      if (value && !Array.isArray(this.value)) {
        this.updateSelectedFilter(this.value ? [this.value] : []);
      } else if (!value && Array.isArray(this.value)) {
        this.updateSelectedFilter(this.value[0] || null);
      }
    },

    updateSelectedFilter(value) {
      this.$emit('input', value);
    },

    updateCondition(condition) {
      this.$emit('update:condition', condition);
    },

    updateFilters(filters, value) {
      this.$emit('update:filters', filters, value);
    },

    createFilter(filter) {
      this.updateFilters([...this.filters, filter]);
    },

    updateFilter(filter, index) {
      const oldFilter = this.filters[index];
      let newValue = this.value;

      if (this.isMultiple) {
        newValue = this.value.map((selectedFilter) => {
          if (selectedFilter.filter === oldFilter.filter) {
            return filter;
          }

          return selectedFilter;
        });
      } else if (this.value && this.value.filter === oldFilter.filter) {
        newValue = filter;
      }

      const newFilters = this.filters.map((selectedFilter, i) => {
        if (i === index) {
          return filter;
        }

        return selectedFilter;
      });

      this.updateFilters(newFilters, newValue);
    },

    deleteFilter(index) {
      const oldFilter = this.filters[index];
      let newValue = this.value;

      if (this.isMultiple) {
        newValue = this.value.filter(selectedFilter => selectedFilter.filter !== oldFilter.filter);
      } else if (this.value && this.value.filter === oldFilter.filter) {
        newValue = null;
      }

      const newFilters = this.filters.filter((filter, i) => i !== index);

      this.updateFilters(newFilters, newValue);
    },

    showFiltersListModal() {
      this.showModal({
        name: MODALS.filtersList,
        config: {
          filters: this.filters,
          hasAccessToAddFilter: this.hasAccessToAddFilter,
          hasAccessToEditFilter: this.hasAccessToEditFilter,
          actions: {
            create: this.createFilter,
            update: this.updateFilter,
            delete: this.deleteFilter,
          },
        },
      });
    },
  },
};
</script>
