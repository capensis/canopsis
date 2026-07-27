<template>
  <c-patterns-field
    :value="patternsField"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
    :with-alarm="isAlarmType"
    :entity-counters-type="!isAlarmType"
    some-required
    with-entity
    @update:value="patternsField = $event"
  />
</template>

<script>
import { useFormModel } from '@/hooks/form';
import { useValidationHeader } from '@/hooks/validator/validation-header';
import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';

export default {
  props: {
    isAlarmType: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const patternsField = useFormModel('patterns');
    const { fetchLinkRulePatternFields } = usePatternsFields();
    const { hasAnyError } = useValidationHeader();

    const { pending, alarmAttributes, entityAttributes } = usePatternsFieldsFetching(
      fetchLinkRulePatternFields,
    );

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs
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
