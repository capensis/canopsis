<template>
  <c-patterns-field
    :value="form"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :readonly="readonly"
    :pending="pending"
    class="mt-2"
    with-alarm
    with-entity
    some-required
    @update:value="$emit('update:form', $event)"
  />
</template>

<script>
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
  },

  emits: ['update:form'],

  setup(props) {
    const { fetchAlarmTagPatternFields } = usePatternsFields();

    const { pending, alarmAttributes, entityAttributes } = usePatternsFieldsFetching(
      fetchAlarmTagPatternFields,
      props.readonly,
    );

    return {
      pending,
      alarmAttributes,
      entityAttributes,
    };
  },
};
</script>
