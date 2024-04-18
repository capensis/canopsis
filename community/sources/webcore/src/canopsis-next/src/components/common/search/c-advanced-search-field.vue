<template>
  <advanced-search-field
    :value="internalValue"
    :fields="fields"
    :conditions="conditions"
    :initial-internal-search="internalSearch"
    @input="updateInternalValue"
  />
</template>

<script>
import { ref, watch } from 'vue';

import { ADVANCED_SEARCH_CONDITIONS } from '@/constants';

import { advancedSearchArrayToString, advancedSearchStringToArray } from '@/helpers/search/advanced-search';

import { useModelField } from '@/hooks/form/model-field';

import AdvancedSearchField from './partials/advanced-search-field.vue';

export default {
  components: { AdvancedSearchField },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    fields: {
      type: Array,
      default: () => [],
    },
    conditions: {
      type: Array,
      default: () => Object.values(ADVANCED_SEARCH_CONDITIONS),
    },
  },
  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);

    const internalStringValue = ref(props.value);

    const updateModelWithInternalValue = (value) => {
      internalStringValue.value = value;
      updateModel(value);
    };

    const {
      value: initialValue,
      internalSearch: initialInternalSearch,
    } = advancedSearchStringToArray(props.value, props.fields);
    const internalValue = ref(initialValue);
    const internalSearch = ref(initialInternalSearch);

    const clearInternalValue = () => {
      internalValue.value = [];
    };

    const clear = () => {
      clearInternalValue();
      updateModelWithInternalValue('');
    };

    const updateInternalValue = (value) => {
      internalValue.value = value;

      updateModelWithInternalValue(advancedSearchArrayToString(value));
    };

    watch(() => props.value, (value) => {
      if (internalStringValue.value === value) {
        return;
      }

      const {
        value: watchInitialValue,
        internalSearch: watchInitialInternalSearch,
      } = advancedSearchStringToArray(value, props.fields);

      internalStringValue.value = value;
      internalSearch.value = watchInitialInternalSearch;
      internalValue.value = watchInitialValue;
    });

    return {
      internalValue,
      internalSearch,

      clear,
      updateInternalValue,
    };
  },
};
</script>
