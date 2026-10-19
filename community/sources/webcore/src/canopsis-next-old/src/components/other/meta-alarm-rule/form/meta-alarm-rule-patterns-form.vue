<template>
  <c-patterns-field
    v-field="form"
    :readonly="readonly"
    :alarm-attributes="alarmAttributes"
    :entity-attributes="entityAttributes"
    :pending="pending"
    :with-total-entity="withTotalEntity"
    :some-required="someRequired"
    :counter-method="counterMethod"
    with-alarm
    with-entity
  >
    <template #additional-counters="{ counters }">
      <v-layout v-if="counters?.threshold_summary" class="gap-2">
        <strong>{{ $t('common.summary') }}:</strong>
        <span class="pre-line">{{ counters.threshold_summary }}</span>
      </v-layout>
    </template>
  </c-patterns-field>
</template>

<script>
import { useMetaAlarmRule } from '@/hooks/store/modules/meta-alarm-rule';
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
