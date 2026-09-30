<template>
  <c-select
    v-model="model"
    :items="preparedFilters"
    :label="label"
    :disabled="disabled"
    :item-title="itemText"
    :item-value="itemValue"
    :multiple="isMultiple"
    :show-count="2"
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

    <template v-if="lockedItems.length" #prepend-inner>
      <v-chip
        v-for="lockedItem in lockedItems"
        :key="getItemValue(lockedItem)"
        class="filter-selector__locked-chip"
        size="small"
      >
        <c-tooltip activator="parent" location="top">
          {{ $t('settings.lockedFilter') }}
        </c-tooltip>
        <span>{{ getItemText(lockedItem) }}</span>
        <v-icon class="ml-2" size="small">
          lock
        </v-icon>
      </v-chip>
    </template>

    <template #selection="{ item, index }">
      <v-chip
        :closable="isChipRemovable"
        size="small"
        @click:close="removeFilter(index)"
      >
        {{ getItemText(item) }}
      </v-chip>
    </template>

    <template #item="{ props, item }">
      <v-list-item v-bind="props" :disabled="isFilterItemDisabled(item)">
        <template v-if="isMultiple" #prepend>
          <v-checkbox-btn :model-value="isFilterSelected(item)" color="primary" />
        </template>

        <template #title>
          <span class="filter-selector__item-title">
          <span class="text-truncate">{{ getItemText(item) }}</span>
          <v-icon
            :color="isFilterSelected(item) ? 'primary' : ''"
            class="ml-2 flex-shrink-0"
            size="small"
          >
            {{ getItemIcon(item) }}
          </v-icon>
          </span>
        </template>
      </v-list-item>
    </template>
  </c-select>
</template>

<script>
import { isArray } from 'lodash';
import { computed } from 'vue';

import { useModelField } from '@/hooks/form/model-field';

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

  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);

    const model = computed({
      get() {
        return props.modelValue;
      },
      set(value) {
        updateModel(value);
      },
    });

    const isMultiple = computed({
      set(isMultiple) {
        if (isMultiple) {
          model.value = props.modelValue ? [props.modelValue] : [];
        } else {
          model.value = props.modelValue?.length ? props.modelValue[0] : undefined;
        }
      },
      get() {
        return isArray(props.modelValue);
      },
    });

    const selectedValues = computed(() => (
      isArray(props.modelValue) ? props.modelValue : [props.modelValue].filter(Boolean)
    ));

    const isOneValueSelected = computed(() => selectedValues.value.length === 1);

    const isChipRemovable = computed(() => {
      if (props.clearable) {
        return true;
      }

      return isMultiple.value ? !isOneValueSelected.value : false;
    });

    const getItemText = item => item?.[props.itemText];
    const getItemValue = item => item?.[props.itemValue];
    const getItemIcon = item => (item?.is_user_preference ? 'person' : 'lock');

    const isLockedFilter = (filter) => {
      const value = getItemValue(filter);

      return isMultiple.value ? props.lockedValue?.includes(value) : props.lockedValue === value;
    };

    const lockedItems = computed(() => props.lockedFilters.filter(isLockedFilter));

    const preparedFilters = computed(() => {
      const filters = [...props.filters];

      if (!props.lockedFilters.length) {
        return filters;
      }

      return [
        ...filters,
        ...props.lockedFilters.map(filter => ({
          ...filter,
          props: {
            disabled: isLockedFilter(filter),
          },
        })),
      ];
    });

    const isFilterSelected = filter => selectedValues.value.includes(getItemValue(filter));

    const isFilterItemDisabled = (filter) => {
      if (isLockedFilter(filter)) {
        return true;
      }

      if (props.clearable) {
        return false;
      }

      const value = getItemValue(filter);

      if (isMultiple.value) {
        return isOneValueSelected.value && selectedValues.value.includes(value);
      }

      return props.modelValue === value;
    };

    const removeFilter = (index) => {
      if (isMultiple.value) {
        model.value = selectedValues.value.filter((value, valueIndex) => valueIndex !== index);
      } else {
        model.value = '';
      }
    };

    return {
      model,
      isMultiple,
      isChipRemovable,
      lockedItems,
      preparedFilters,
      getItemText,
      getItemValue,
      getItemIcon,
      isFilterSelected,
      isFilterItemDisabled,
      removeFilter,
    };
  },
};
</script>

<style lang="scss">
.filter-selector__locked-chip {
  margin-right: 4px;
}
</style>
