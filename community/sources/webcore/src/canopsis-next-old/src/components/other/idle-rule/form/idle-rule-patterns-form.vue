<template>
  <c-patterns-field
    v-field="form"
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
