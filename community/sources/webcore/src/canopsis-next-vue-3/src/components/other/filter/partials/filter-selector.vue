<template>
  <c-select
    v-model="model"
    :items="preparedFilters"
    :label="label"
    :disabled="disabled"
    :item-title="itemText"
    :item-value="itemValue"
    :multiple="isMultiple"
    :hide-details="hideDetails"
    :clearable="clearable"
    class="filter-selector mt-0"
  >
    <template v-if="!hideMultiply" #prepend-item>
      <c-enabled-field
        v-model="isMultiple"
        :label="$t('filter.selector.fields.mixFilters')"
        class="mx-3"
        hide-details
      />
      <v-divider class="mt-3" />
    </template>

    <template #selection="{ item, index }">
      <v-chip
        v-if="index < 2"
        :closable="isChipRemovable"
        size="small"
        @click:close="removeFilter(index)"
      >
        {{ getItemText(item.raw ?? item) }}
      </v-chip>
    </template>

    <template #item="{ props, item }">
      <v-list-item
        v-bind="props"
        :disabled="isFilterItemDisabled(item.raw)"
      >
        <template v-if="isMultiple" #prepend>
          <v-checkbox-btn
            :model-value="isFilterSelected(item.raw)"
            color="primary"
          />
        </template>

        <v-list-item-title>
          <span>{{ getItemText(item.raw) }}</span>
          <v-icon
            :color="isFilterSelected(item.raw) ? 'primary' : ''"
            class="ml-2"
            size="small"
          >
            {{ getItemIcon(item.raw) }}
          </v-icon>
        </v-list-item-title>
      </v-list-item>
    </template>
  </c-select>
</template>

<script>
import { isArray } from 'lodash';

export default {
  props: {
    modelValue: {
      type: [String, Array],
      default: null,
    },
    lockedValue: {
      type: [String, Array],
      default: null,
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

  emits: ['update:modelValue'],

  computed: {
    model: {
      get() {
        return this.modelValue;
      },
      set(value) {
        this.$emit('update:modelValue', value);
      },
    },

    isMultiple: {
      set(isMultiple) {
        if (isMultiple) {
          this.model = this.modelValue ? [this.modelValue] : [];
        } else {
          this.model = this.modelValue?.length ? this.modelValue[0] : undefined;
        }
      },
      get() {
        return isArray(this.modelValue);
      },
    },

    selectedValues() {
      return isArray(this.modelValue) ? this.modelValue : [this.modelValue].filter(Boolean);
    },

    isOneValueSelected() {
      return this.selectedValues.length === 1;
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

      return [
        ...preparedFilters,
        ...this.lockedFilters.map(filter => ({
          ...filter,
          props: {
            disabled: this.isLockedFilter(filter),
          },
        })),
      ];
    },
  },

  methods: {
    getItemText(item) {
      return item?.[this.itemText];
    },

    getItemValue(item) {
      return item?.[this.itemValue];
    },

    getItemIcon(item) {
      return item?.is_user_preference ? 'person' : 'lock';
    },

    isFilterSelected(filter) {
      return this.selectedValues.includes(this.getItemValue(filter));
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
        return this.isOneValueSelected && this.selectedValues.includes(value);
      }

      return this.modelValue === value;
    },

    isLockedFilter(filter) {
      const value = this.getItemValue(filter);

      return this.isMultiple
        ? this.lockedValue?.includes(value)
        : this.lockedValue === value;
    },

    removeFilter(index) {
      if (this.isMultiple) {
        this.model = this.selectedValues.filter((value, valueIndex) => valueIndex !== index);
      } else {
        this.model = '';
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
