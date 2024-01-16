<template>
  <v-select
    class="filter-selector mt-0"
    v-field="value"
    :items="preparedFilters"
    :label="label"
    :disabled="disabled"
    :always-dirty="!!lockedItems.length"
    :item-text="itemText"
    :item-value="itemValue"
    :item-disabled="isFilterItemDisabled"
    :multiple="isMultiple"
    :hide-details="hideDetails"
  >
    <template
      v-if="!hideMultiply"
      #prepend-item=""
    >
      <c-enabled-field
        class="mx-3"
        v-model="isMultiple"
        :label="$t('filter.selector.fields.mixFilters')"
        hide-details
      />
      <v-divider class="mt-3" />
    </template>
    <template #selections="{ items }">
      <v-tooltip
        v-for="lockedItem in lockedItems"
        :key="getItemValue(lockedItem)"
        top
      >
        <template #activator="{ on }">
          <v-chip
            v-on="on"
            small
          >
            <span>{{ getItemText(lockedItem) }}</span>
            <v-icon
              class="ml-2"
              small
            >
              lock
            </v-icon>
          </v-chip>
        </template>
        <span>{{ $t('settings.lockedFilter') }}</span>
      </v-tooltip>
      <v-chip
        v-for="(item, index) in items"
        :key="getItemValue(item)"
        :close="isChipRemovable"
        small
        @click:close="removeFilter(index)"
      >
        {{ getItemText(item) }}
      </v-chip>
    </template>
    <template #item="{ parent, item, attrs, on }">
      <v-list-item
        v-bind="attrs"
        v-on="on"
      >
        <v-list-item-action v-if="isMultiple">
          <v-checkbox
            :input-value="item.active || attrs.value"
            :color="parent.color"
            :disabled="attrs.disabled"
          />
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>
            <span>{{ item.title }}</span>
            <v-icon
              class="ml-2"
              :color="attrs.value ? parent.color : ''"
              small
            >
              {{ getItemIcon(item) }}
            </v-icon>
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </template>
  </v-select>
</template>

<script>
import { isArray } from 'lodash';

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
      set(isMultiple) {
        if (isMultiple) {
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
      return item.is_user_preference ? 'person' : 'lock';
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
