<template>
  <c-patterns-field
    v-field="form"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :readonly="readonly"
    :pending="pending"
    class="mt-2"
    with-alarm
    with-entity
    some-required
  />
</template>

<script>
import { useValidationHeader } from '@/hooks/validator/validation-header';
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';

export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
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
  setup() {
    const { fetchAlarmTagPatternFields } = usePatternsFields();
    const { hasAnyError } = useValidationHeader();

    const {
      pending,
      alarmAttributes,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchAlarmTagPatternFields);

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
