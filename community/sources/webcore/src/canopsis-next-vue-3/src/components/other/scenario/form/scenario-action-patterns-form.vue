<template>
  <c-patterns-field
    v-model:value="patternsField"
    :name="name"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
    some-required
    with-alarm
    with-entity
  />
</template>

<script>
import { computed } from 'vue';

import { useModelField } from '@/hooks/form/model-field';
import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';
import { useValidationHeader } from '@/hooks/validator/validation-header';

export default {
  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'patterns',
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { fetchScenarioPatternFields } = usePatternsFields();
    const { hasAnyError } = useValidationHeader();
    const { updateModel } = useModelField(props, emit);
    const patternsField = computed({
      get: () => props.modelValue,
      set: value => updateModel(value),
    });

    const { pending, alarmAttributes, entityAttributes } = usePatternsFieldsFetching(
      fetchScenarioPatternFields,
    );

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs.
       */
      hasAnyError,
      patternsField,
      pending,
      alarmAttributes,
      entityAttributes,
    };
  },
};
</script>
