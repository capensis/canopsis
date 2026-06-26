<template>
  <c-patterns-field
    :value="form"
    :with-alarm="!isEntityType"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
    :readonly="readonly"
    :entity-counters-type="isEntityType"
    some-required
    with-entity
    @update:value="$emit('update:form', $event)"
  />
</template>

<script>
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';

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

  emits: ['update:form'],

  setup(props) {
    const { fetchIdleRulePatternFields } = usePatternsFields();

    const {
      pending,
      alarmAttributes,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchIdleRulePatternFields, props.readonly);

    return {
      pending,
      alarmAttributes,
      entityAttributes,
    };
  },
};
</script>
