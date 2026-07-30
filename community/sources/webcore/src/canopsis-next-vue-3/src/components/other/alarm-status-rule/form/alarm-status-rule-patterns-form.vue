<template>
  <c-patterns-field
    v-model:value="patternsField"
    :readonly="readonly"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
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

  emits: ['update:form'],

  setup(props, { emit }) {
    const patternsField = useFormModel('patterns');

    const { fetchFlappingRulePatternFields, fetchResolveRulePatternFields } = usePatternsFields();

    const fetchAction = computed(() =>
      props.flapping ? fetchFlappingRulePatternFields : fetchResolveRulePatternFields,
    );

    const { pending, alarmAttributes, entityAttributes } = usePatternsFieldsFetching(
      fetchAction,
      props.readonly,
    );

    const updateForm = value => emit('update:form', value);

    return {
      patternsField,
      pending,
      alarmAttributes,
      entityAttributes,
      updateForm,
    };
  },
};
</script>
