<template>
  <c-patterns-field
    v-model:value="patternsField"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :readonly="readonly"
    :pending="pending"
    some-required
    with-pbehavior
    with-alarm
    with-entity
  />
</template>

<script>
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';
import { useFormModel } from '@/hooks/form';
import { useValidationHeader } from '@/hooks/validator/validation-header';

export default {
  props: {
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { fetchDeclareTicketRulePatternFields } = usePatternsFields();
    const patternsField = useFormModel('patterns');
    const { hasAnyError } = useValidationHeader();

    const {
      pending,
      alarmAttributes,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchDeclareTicketRulePatternFields, props.readonly);

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
