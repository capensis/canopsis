<template>
  <advanced-search-field
    :value="internalValue"
    :fields="fields"
    :conditions="conditions"
    :error-messages="errorMessages"
    @input="updateInternalValue"
  />
</template>

<script>
import { ref, computed, watch } from 'vue';

import { ADVANCED_SEARCH_CONDITIONS, ADVANCED_SEARCH_ITEM_TYPES, ALARM_ADVANCED_SEARCH_FIELDS } from '@/constants';

import { advancedSearchStringToArray, advancedSearchArrayToString } from '@/helpers/search/advanced-search';

import { useI18n } from '@/hooks/i18n';
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
    columns: {
      type: Array,
      default: () => [],
    },
    conditions: {
      type: Array,
      default: () => Object.values(ADVANCED_SEARCH_CONDITIONS),
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);

    const fields = computed(() => (
      props.columns
        .filter(({ value }) => ALARM_ADVANCED_SEARCH_FIELDS.includes(value))
        .map(({ text }) => ({ text, value: text, type: ADVANCED_SEARCH_ITEM_TYPES.field }))
    ));

    const errorMessages = ref([]);
    const internalStringValue = ref(props.value);

    const updateModelWithInternalValue = (value) => {
      internalStringValue.value = value;
      updateModel(value);
    };

    const getParsedValue = (value) => {
      try {
        return advancedSearchStringToArray(value);
      } catch (err) {
        console.error(err);
        errorMessages.value = [t('advancedSearch.validationError')];

        return [];
      }
    };

    const internalValue = ref(getParsedValue(props.value));

    const clearInternalValue = () => {
      internalValue.value = [];
    };

    const clearErrorMessages = () => {
      errorMessages.value = [];
    };

    const clear = () => {
      clearInternalValue();
      clearErrorMessages();
      updateModelWithInternalValue('');
    };

    const updateInternalValue = (value) => {
      internalValue.value = value;

      clearErrorMessages();
      updateModelWithInternalValue(advancedSearchArrayToString(value));
    };

    watch(() => props.value, (value) => {
      if (internalStringValue.value === value) {
        return;
      }

      internalStringValue.value = value;
      internalValue.value = getParsedValue(value);
    });

    return {
      fields,
      internalValue,
      errorMessages,

      clear,
      clearErrorMessages,
      updateInternalValue,
    };
  },
};
</script>
