<template>
  <c-patterns-field
    v-model:value="patternsField"
    :readonly="readonly"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
    :expanded-alarm="readonly"
    :expanded-entity="readonly"
    name="patterns"
    with-alarm
    with-entity
    some-required
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
    readonly: {
      type: Boolean,
      default: false,
    },
    flapping: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const patternsField = props.readonly
     ? computed(() => props.form)
     : useFormModel('patterns');

    const { fetchFlappingRulePatternFields, fetchResolveRulePatternFields } = usePatternsFields();

    const fetchAction = computed(() => props.flapping ? fetchFlappingRulePatternFields : fetchResolveRulePatternFields,
    );

    const { pending, alarmAttributes, entityAttributes } = usePatternsFieldsFetching(
      fetchAction,
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
