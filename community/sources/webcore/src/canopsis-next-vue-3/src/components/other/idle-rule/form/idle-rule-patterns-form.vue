<template>
  <c-patterns-field
    v-model:value="patternsField"
    :with-alarm="!isEntityType"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
    :readonly="readonly"
    :entity-counters-type="isEntityType"
    some-required
    with-entity
  />
</template>

<script>
import { computed } from 'vue';

import { useFormModel } from '@/hooks/form';
import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';

export default {
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isEntityType: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const patternsField = props.readonly
     ? computed(() => props.form)
     : useFormModel('patterns');

    const { fetchIdleRulePatternFields } = usePatternsFields();

    const { pending, alarmAttributes, entityAttributes } = usePatternsFieldsFetching(
      fetchIdleRulePatternFields,
      props.readonly,
    );

    return {
      patternsField,
      pending,
      alarmAttributes,
      entityAttributes,
    };
  },
};
</script>
