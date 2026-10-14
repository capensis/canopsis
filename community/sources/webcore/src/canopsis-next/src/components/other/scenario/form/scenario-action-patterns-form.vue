<template>
  <c-patterns-field
    :value="modelValue"
    :name="name"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
    some-required
    with-alarm
    with-entity
    @update:value="updateModel"
  />
</template>

<script>
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

    const { pending, alarmAttributes, entityAttributes } = usePatternsFieldsFetching(
      fetchScenarioPatternFields,
    );

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs.
       */
      hasAnyError,
      updateModel,
      pending,
      alarmAttributes,
      entityAttributes,
    };
  },
};
</script>
