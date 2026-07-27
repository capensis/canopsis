<template>
  <c-patterns-field
    :value="form"
    :readonly="readonly"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
    :with-total-entity="withTotalEntity"
    :some-required="someRequired"
    :counter-method="counterMethod"
    with-alarm
    with-entity
    @update:value="$emit('update:form', $event)"
  >
    <template #additional-counters="{ counters }">
      <c-layout v-if="counters?.threshold_summary" gap="2">
        <strong>{{ $t('common.summary') }}:</strong>
        <span class="pre-line">{{ counters.threshold_summary }}</span>
      </c-layout>
    </template>
  </c-patterns-field>
</template>

<script>
import { useMetaAlarmRule } from '@/hooks/store/modules/meta-alarm-rule';
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
    withTotalEntity: {
      type: Boolean,
      default: false,
    },
    someRequired: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:form'],

  setup(props) {
    const { fetchMetaalarmrulePatternFields } = usePatternsFields();

    const { pending, alarmAttributes, entityAttributes } = usePatternsFieldsFetching(
      fetchMetaalarmrulePatternFields,
      props.readonly,
    );

    const { checkPatternsAlarmsCount } = useMetaAlarmRule();

    return {
      pending,
      alarmAttributes,
      entityAttributes,
      counterMethod: checkPatternsAlarmsCount,
    };
  },
};
</script>
