<template>
  <c-patterns-field
    v-field="form"
    :readonly="readonly"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
    with-alarm
    with-entity
    some-required
  />
</template>

<script>
import { computed } from 'vue';

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
    flapping: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { fetchFlappingRulePatternFields, fetchResolveRulePatternFields } = usePatternsFields();

    const fetchAction = computed(() => (
      props.flapping
        ? fetchFlappingRulePatternFields
        : fetchResolveRulePatternFields
    ));

    const {
      pending,
      alarmAttributes,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchAction, props.readonly);

    return {
      pending,
      alarmAttributes,
      entityAttributes,
    };
  },
};
</script>
