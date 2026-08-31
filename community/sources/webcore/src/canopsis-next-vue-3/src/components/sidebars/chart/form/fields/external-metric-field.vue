<template>
  <c-lazy-search-field
    ref="field"
    :value="selectedItem"
    :items="preparedItems"
    :label="label"
    :name="name"
    :loading="wholePending"
    :has-more="hasMoreItems"
    :hide-no-data="true"
    :menu="menu"
    :menu-icon="addable ? 'arrow_drop_down' : undefined"
    :required="required"
    :search="searchValue"
    item-text="text"
    item-value="value"
    return-object
    @keyup.enter.capture="handleEnter"
    @update:menu="menu = $event"
    @update:value="handleSelectedValue"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="handleSearchUpdate"
  >
    <template v-if="!addable" #selection="{ item }">
      <v-icon class="mr-2">
        language
      </v-icon>
      <span>{{ item?.value ?? item }}</span>
    </template>

    <template v-if="!addable" #icon>
      <v-icon class="mr-2">
        language
      </v-icon>
    </template>
  </c-lazy-search-field>
</template>

<script>
import { isString } from 'lodash';
import { computed, nextTick, ref, toRef } from 'vue';

import { useLazySearch } from '@/hooks/form/lazy-search';
import { useMetrics } from '@/hooks/store/modules/metrics';

export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    addable: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'metric',
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const menu = ref(false);
    const searchValue = ref('');
    const field = ref(null);
    const { fetchExternalMetricsWithoutStore } = useMetrics();
    const fetchExternalMetrics = async (params) => {
      const response = await fetchExternalMetricsWithoutStore(params);

      return {
        ...response,
        data: response.data.map(metric => isString(metric) ? { value: metric } : metric,
        ),
      };
    };
    const emitModelValue = (event, value) => {
      if (event === 'update:selected-items') {
        return;
      }

      emit(event === 'update:value' ? 'update:modelValue' : event, value?.value ?? value);
    };
    const lazySearch = useLazySearch(
      {
        value: toRef(props, 'modelValue'),
        idKey: 'value',
        idParamsKey: 'ids',
        fetchHandler: fetchExternalMetrics,
        addable: props.addable,
        attachValue: true,
      },
      emitModelValue,
    );
    const handleSearchUpdate = (value) => {
      searchValue.value = value;
      lazySearch.updateSearch(value);
    };
    const handleSelectedValue = (value) => {
      if (isString(value)) {
        emit('update:modelValue', value);
        closeMenuAndBlur();

        return;
      }

      lazySearch.changeSelectedItems(value);
    };
    const closeMenuAndBlur = () => {
      menu.value = false;
      field.value?.$el?.querySelector('input')?.blur();

      nextTick(() => {
        menu.value = false;
        field.value?.$el?.querySelector('input')?.blur();
      });
    };
    const handleEnter = (event) => {
      const enteredValue =
        field.value?.$el?.querySelector('input')?.value ?? event.target?.value ?? searchValue.value;
      searchValue.value = enteredValue;

      const normalizedSearch = enteredValue?.trim().toLowerCase() || '';

      if (!normalizedSearch) {
        return;
      }

      const matchedMetric = lazySearch.items.value.find(
        ({ value }) => value?.trim().toLowerCase() === normalizedSearch,
      );

      event.preventDefault();
      event.stopPropagation();
      emit('update:modelValue', matchedMetric?.value ?? enteredValue.trim());
      closeMenuAndBlur();
    };
    const preparedItems = computed(() => lazySearch.items.value.map(({ value }) => ({ value, text: value })),
    );

    return {
      ...lazySearch,
      field,
      menu,
      preparedItems,
      searchValue,
      handleEnter,
      handleSearchUpdate,
      handleSelectedValue,
    };
  },
};
</script>
