<template lang="pug">
  v-layout(align-center, row, wrap)
    v-flex(
      data-test="selectFilters",
      v-show="!hideSelect",
      v-bind="flexProps.select"
    )
      v-select(
        :value="value",
        :items="preparedFilters",
        :label="label",
        :itemText="itemText",
        :itemValue="itemValue",
        :multiple="isMultiple",
        :disabled="!hasAccessToListFilters && !hasAccessToUserFilter",
        return-object,
        clearable,
        @input="updateSelectedFilter"
      )
        template(slot="prepend-item")
          v-layout.pl-3
            v-flex(v-show="!hideSelect", v-bind="flexProps.switch")
              v-switch(
                data-test="mixFilters",
                color="primary",
                :label="$t('filterSelector.fields.mixFilters')",
                :input-value="isMultiple",
                :disabled="!hasAccessToListFilters && !hasAccessToUserFilter",
                hide-details,
                @change="updateIsMultipleFlag"
              )
            v-flex(v-show="!hideSelect && isMultiple", v-bind="flexProps.radio")
              v-radio-group.mb-0(
                :value="condition",
                :disabled="!hasAccessToListFilters && !hasAccessToUserFilter",
                hide-details,
                row,
                @change="updateCondition"
              )
                v-radio(
                  :value="$constants.FILTER_MONGO_OPERATORS.and",
                  data-test="andFilters",
                  label="AND",
                  color="primary"
                )
                v-radio(
                  :value="$constants.FILTER_MONGO_OPERATORS.or",
                  data-test="orFilters",
                  label="OR",
                  color="primary"
                )
          v-divider.mt-3
        template(slot="item", slot-scope="{ parent, item, tile }")
          v-list-tile-action(
            v-if="isMultiple",
            :data-test="`filterOption-${item[itemText]}`",
            @click.stop="parent.$emit('select', item)"
          )
            v-checkbox(:inputValue="tile.props.value", :color="parent.color")
          v-list-tile-content
            v-list-tile-title
              span {{ item[itemText] }}
              v-icon.ml-2(
                v-show="!hideSelectIcon",
                :color="tile.props.value ? parent.color : ''",
                small
              ) {{ item.locked ? 'lock' : 'person' }}
    v-flex(v-if="hasAccessToUserFilter", v-bind="flexProps.list")
      v-tooltip(v-if="!long", bottom)
        v-btn(
          slot="activator",
          data-test="showFiltersListButton",
          icon,
          small,
          @click="showFiltersListModal"
        )
          v-icon filter_list
        span {{ $t('filterSelector.buttons.list') }}
      filters-form(
        v-else,
        :filters="filtersWithSelected",
        :entitiesType="entitiesType",
        @input="updateFilters"
      )
</template>

<script>
import { isEmpty, omit } from 'lodash';

import { ENTITIES_TYPES, MODALS, FILTER_DEFAULT_VALUES } from '@/constants';

import FiltersForm from '@/components/other/filter/form/filters-form.vue';

export default {
  components: { FiltersForm },
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
    entitiesType: {
      type: String,
      default: ENTITIES_TYPES.alarm,
      validator: value => [ENTITIES_TYPES.alarm, ENTITIES_TYPES.entity].includes(value),
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
      const preparedFilters = this.hasAccessToUserFilter ? [...this.filters] : [];
      const preparedLockedFilters = this.lockedFilters.map(filter => ({ ...filter, locked: true }));

      if (preparedFilters.length && preparedLockedFilters.length) {
        return preparedFilters.concat({ divider: true }, preparedLockedFilters);
      } if (preparedFilters.length) {
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
        this.updateSelectedFilter(!isEmpty(this.value) ? [this.value] : []);
      } else if (!checked && isValueArray) {
        this.updateSelectedFilter(!isEmpty(this.value[0]) ? this.value[0] : null);
      }
    },

    updateSelectedFilter(newValue) {
      this.$emit('input', newValue);
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

      return this.$emit('update:filters', filters.map(removeSelectedProperty), newValue);
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
          entitiesType: this.entitiesType,
          action: filters => this.updateFilters(filters),
        },
      });
    },
  },
};
</script>
