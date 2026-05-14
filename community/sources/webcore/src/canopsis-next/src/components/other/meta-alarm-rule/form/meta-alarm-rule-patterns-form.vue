<template>
  <c-patterns-field
    v-field="form"
    :readonly="readonly"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
    :with-total-entity="withTotalEntity"
    :some-required="someRequired"
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
    withTotalEntity: {
      type: Boolean,
      default: false,
    },
    someRequired: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { fetchMetaalarmrulePatternFields } = usePatternsFields();

    const {
      pending,
      alarmAttributes,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchMetaalarmrulePatternFields, props.readonly);

    return {
      pending,
      alarmAttributes,
      entityAttributes,
    };
  },
};
</script>
