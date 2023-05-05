<template lang="pug">
  v-select.filter-selector(
    v-field="value",
    :items="preparedFilters",
    :label="label",
    :disabled="disabled",
    :always-dirty="!!lockedItems.length",
    :item-text="itemText",
    :item-value="itemValue",
    :item-disabled="isFilterItemDisabled",
    :multiple="isMultiple",
    :hide-details="hideDetails"
  )
    template(v-if="!hideMultiply", #prepend-item="")
      c-enabled-field.mx-3(
        v-model="isMultiple",
        :label="$t('filter.selector.fields.mixFilters')",
        hide-details
      )
      v-divider.mt-3

    template(#selections="{ items }")
      v-tooltip(
        v-for="lockedItem in lockedItems",
        :key="getItemValue(lockedItem)",
        top
      )
        template(#activator="{ on }")
          v-chip(v-on="on", small)
            span {{ getItemText(lockedItem) }}
            v-icon.ml-2(small) lock
        span {{ $t('settings.lockedFilter') }}
      v-chip(
        v-for="(item, index) in items",
        :key="getItemValue(item)",
        :close="isChipRemovable",
        small,
        @input="removeFilter(index)"
      ) {{ getItemText(item) }}

    template(#item="{ parent, item, tile }")
      v-list-tile(v-bind="tile.props", v-on="tile.on")
        v-list-tile-action(v-if="isMultiple")
          v-checkbox(
            :input-value="item.active || tile.props.value",
            :color="parent.color",
            :disabled="tile.props.disabled"
          )
        v-list-tile-content
          v-list-tile-title.v-list-badge__tile__title
            span {{ item.title }}
            v-badge(
              :value="isOldPattern(item)",
              color="error",
              overlap
            )
              template(#badge="")
                v-tooltip(top)
                  template(#activator="{ on: badgeTooltipOn }")
                    v-icon(v-on="badgeTooltipOn", color="white") priority_high
                  span {{ $t('pattern.oldPatternTooltip') }}
              v-icon.ml-2(
                :color="tile.props.value ? parent.color : ''",
                small
              ) {{ getItemIcon(item) }}
</template>

<script>
import { isArray } from 'lodash';

import { isOldPattern } from '@/helpers/pattern';

import { formArrayMixin } from '@/mixins/form';

export default {
  mixins: [formArrayMixin],
  props: {
    value: {
      type: [String, Array],
      default: () => null,
    },
    lockedValue: {
      type: [String, Array],
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
      default: '_id',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: true,
    },
    hideMultiply: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isMultiple: {
      set(value) {
        if (value) {
          this.updateModel(this.value ? [this.value] : []);
        } else {
          this.updateModel(this.value.length ? this.value[0] : undefined);
        }
      },
      get() {
        return isArray(this.value);
      },
    },

    isOneValueSelected() {
      return this.value.length === 1;
    },

    isChipRemovable() {
      if (this.clearable) {
        return true;
      }

      return this.isMultiple ? !this.isOneValueSelected : false;
    },

    lockedItems() {
      return this.lockedFilters.filter(this.isLockedFilter);
    },

    preparedFilters() {
      const preparedFilters = [...this.filters];

      if (!this.lockedFilters.length) {
        return preparedFilters;
      }

      if (preparedFilters.length) {
        preparedFilters.push({ divider: true });
      }

      preparedFilters.push(
        ...this.lockedFilters.map(filter => ({
          ...filter,
          active: this.isLockedFilter(filter),
        })),
      );

      return preparedFilters;
    },
  },
  methods: {
    getItemText(item) {
      return item[this.itemText];
    },

    getItemValue(item) {
      return item[this.itemValue];
    },

    getItemIcon(item) {
      return item.is_private ? 'person' : 'lock';
    },

    isOldPattern(filter) {
      return isOldPattern(filter);
    },

    isFilterItemDisabled(filter) {
      if (this.isLockedFilter(filter)) {
        return true;
      }

      if (this.clearable) {
        return false;
      }

      const value = this.getItemValue(filter);

      if (this.isMultiple) {
        return this.isOneValueSelected && this.value.includes(value);
      }

      return this.value === value;
    },

    isLockedFilter(filter) {
      const value = this.getItemValue(filter);

      return this.isMultiple
        ? this.lockedValue?.includes(value)
        : this.lockedValue === value;
    },

    removeFilter(index) {
      if (this.isMultiple) {
        this.removeItemFromArray(index);
      } else {
        this.updateModel('');
      }
    },
  },
};
</script>

<style lang="scss">
.filter-selector {
  max-width: 500px;
}
</style>
