<template>
  <c-patterns-field
    :value="patternsValue"
    :readonly="readonly"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
    with-alarm
    with-entity
    some-required
    @update:value="updatePatterns"
  />
</template>

<script>
import { computed, inject } from 'vue';

import { FORM_CONTEXT_KEY } from '@/constants';

import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';
import { useValidationHeader } from '@/hooks/validator/validation-header';

export default {
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const formContext = inject(FORM_CONTEXT_KEY, null);
    const patternsField = formContext?.model('patterns');
    const { hasAnyError } = useValidationHeader();
    const { fetchDynamicInfosPatternFields } = usePatternsFields();

    const {
      pending,
      alarmAttributes,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchDynamicInfosPatternFields, props.readonly);

    const patternsValue = computed(() => patternsField?.value ?? props.form);

    const updatePatterns = (value) => {
      if (patternsField) {
        patternsField.value = value;
      }

      emit('update:form', value);
    };

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs.
       */
      hasAnyError,
      patternsValue,
      pending,
      alarmAttributes,
      entityAttributes,
      updatePatterns,
    };
  },
};
</script>
