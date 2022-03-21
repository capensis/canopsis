<template lang="pug">
  v-layout(align-center, row, wrap)
    v-flex(v-show="!hideSelect", :xs12="long")
      v-select(
        v-field="value",
        :items="preparedFilters",
        :label="label",
        :item-text="itemText",
        :item-value="itemValue",
        :multiple="isMultiple",
        :disabled="!hasAccessToListFilters && !hasAccessToUserFilter",
        return-object,
        clearable
      )
        template(#prepend-item="")
          v-layout.pl-3
            v-flex(v-show="!hideSelect", :xs6="long")
              c-enabled-field(
                :value="isMultiple",
                :label="$t('filterSelector.fields.mixFilters')",
                :disabled="!hasAccessToListFilters && !hasAccessToUserFilter",
                hide-details,
                @input="updateIsMultipleFlag"
              )
            v-flex(v-show="!hideSelect && isMultiple", :xs6="long")
              c-operator-field(:value="condition", @input="updateCondition")
          v-divider.mt-3

        template(#item="{ parent, item, tile }")
          v-list-tile-action(v-if="isMultiple", @click.stop="parent.$emit('select', item)")
            v-checkbox(:input-value="tile.props.value", :color="parent.color")
          v-list-tile-content
            v-list-tile-title
              span {{ item[itemText] }}
              v-icon.ml-2(
                v-show="!hideSelectIcon",
                :color="tile.props.value ? parent.color : ''",
                small
              ) {{ item.locked ? 'lock' : 'person' }}

    v-flex(v-if="hasAccessToUserFilter", :xs12="long")
      c-action-btn(
        v-if="!long",
        :tooltip="$t('filterSelector.buttons.list')",
        icon="filter_list",
        small,
        @click="showFiltersListModal"
      )
      filters-list-form(
        v-else,
        :filters="filtersWithSelected",
        addable,
        editable,
        @input="updateFilters"
      )
</template>

<script>
import { isEmpty, omit } from 'lodash';

import { MODALS, FILTER_DEFAULT_VALUES } from '@/constants';

import { formMixin } from '@/mixins/form';

import FiltersListForm from '@/components/forms/filters/filters-list-form.vue';

export default {
  components: { FiltersListForm },
  mixins: [formMixin],
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
      default: 'title',
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
    hasAccessToListFilters: {
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
    hasAccessToUserFilter: {
      type: Boolean,
      default: true,
    },
    entity: {
      type: Boolean,
      default: false,
    },
    pbehavior: {
      type: Boolean,
      default: false,
    },
    alarm: {
      type: Boolean,
      default() {
        return !this.entity && !this.pbehavior;
      },
    },
  },
  computed: {
    isMultiple() {
      return Array.isArray(this.value);
    },

    preparedFilters() {
      const preparedFilters = this.hasAccessToUserFilter ? [...this.filters] : [];
      const preparedLockedFilters = this.lockedFilters.map(filter => ({ ...filter, locked: true }));

      if (preparedFilters.length && preparedLockedFilters.length) {
        return preparedFilters.concat({ divider: true }, preparedLockedFilters);
      }

      if (preparedFilters.length) {
        return preparedFilters;
      }

      return preparedLockedFilters;
    },

    filtersWithSelected() {
      return this.filters.map((filter) => {
        const selected = this.isMultiple
          ? this.value.some(currentFilter => this.isFilterEqual(filter, currentFilter))
          : !!this.value && this.isFilterEqual(filter, this.value);

        return { ...filter, selected };
      });
    },
  },
  methods: {
    updateIsMultipleFlag(checked) {
      const isValueArray = Array.isArray(this.value);

      if (checked && !isValueArray) {
        this.updateModel(!isEmpty(this.value) ? [this.value] : []);
      } else if (!checked && isValueArray) {
        this.updateModel(!isEmpty(this.value[0]) ? this.value[0] : null);
      }
    },

    updateCondition(newCondition) {
      this.$emit('update:condition', newCondition);
    },

    updateFilters(filters) {
      const removeSelectedProperty = filter => omit(filter, 'selected');

      const selectedFilters = filters.reduce((acc, filter) => {
        if (filter.selected) {
          acc.push(removeSelectedProperty(filter));
        }

        return acc;
      }, []);

      const newValue = this.isMultiple ? selectedFilters : selectedFilters[0];

      this.$emit('update:filters', filters.map(removeSelectedProperty), newValue);
    },

    isFilterEqual(firstFilter, secondFilter) {
      return firstFilter.title === secondFilter.title && firstFilter.filter === secondFilter.filter;
    },

    showFiltersListModal() {
      this.$modals.show({
        name: MODALS.filtersList,
        config: {
          filters: this.filtersWithSelected,
          hasAccessToAddFilter: this.hasAccessToUserFilter,
          hasAccessToEditFilter: this.hasAccessToUserFilter,
          alarm: this.alarm,
          entity: this.entity,
          pbehavior: this.pbehavior,
          action: filters => this.updateFilters(filters),
        },
      });
    },
  },
};
</script>
