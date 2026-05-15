<template>
  <c-patterns-field
    v-field="form"
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
  },
  setup(props) {
    const { fetchDeclareTicketRulePatternFields } = usePatternsFields();

    const {
      pending,
      alarmAttributes,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchDeclareTicketRulePatternFields, props.readonly);

    return {
      pending,
      alarmAttributes,
      entityAttributes,
    };
  },
};
</script>
