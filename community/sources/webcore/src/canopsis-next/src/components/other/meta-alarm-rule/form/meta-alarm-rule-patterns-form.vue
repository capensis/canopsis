<template>
  <c-patterns-field
    v-field="form"
    :readonly="readonly"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
    :with-total-entity="withTotalEntity"
    :some-required="someRequired"
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
    withTotalEntity: {
      type: Boolean,
      default: false,
    },
    someRequired: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { fetchMetaalarmrulePatternFields } = usePatternsFields();
    const { hasAnyError } = useValidationHeader();

    const {
      pending,
      alarmAttributes,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchMetaalarmrulePatternFields);

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
