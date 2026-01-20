<template>
  <c-patterns-field
    v-field="patterns"
    :name="name"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
    some-required
    with-alarm
    with-entity
    both-counters
  />
</template>

<script>
import { useValidationHeader } from '@/hooks/validator/validation-header';
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';

export default {
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'patterns',
    },
  },
  setup() {
    const { fetchScenarioPatternFields } = usePatternsFields();
    const { hasAnyError } = useValidationHeader();

    const {
      pending,
      alarmAttributes,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchScenarioPatternFields);

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs
       */
      hasAnyError,
      pending,
      alarmAttributes,
      entityAttributes,
    };
  },
};
</script>
