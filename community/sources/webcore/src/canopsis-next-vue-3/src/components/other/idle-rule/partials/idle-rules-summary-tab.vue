<template>
  <c-layout column>
    <idle-rules-summary-row :label="t('common.id')" :value="idleRule._id" />
    <idle-rules-summary-row :label="t('common.name')" :value="idleRule.name" />
    <idle-rules-summary-row :label="t('common.description')" :value="idleRule.description" />
    <idle-rules-summary-row
      :label="t('common.type')"
      :value="t(`idleRules.types.${idleRule.type}`)"
    />
    <idle-rules-summary-row
      :label="t('common.created')"
      :value="convertDateToString(idleRule.created)"
    />
    <idle-rules-summary-row
      :label="t('common.updated')"
      :value="convertDateToString(idleRule.updated)"
    />
    <idle-rules-summary-row
      :label="t('idleRules.timeAwaiting')"
      :value="convertDurationToString(idleRule.duration)"
    />
    <idle-rules-summary-row :label="t('common.priority')" :value="idleRule.priority" />
    <idle-rules-summary-row
      v-if="disableDuringPeriods"
      :label="t('common.disableDuringPeriods')"
      :value="disableDuringPeriods"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';
import { convertDurationToString } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';

import IdleRulesSummaryRow from './idle-rules-summary-row.vue';

export default {
  components: { IdleRulesSummaryRow },

  props: {
    idleRule: {
      type: Object,
      default: () => ({}),
    },
  },

  setup(props) {
    const { t } = useI18n();
    const disableDuringPeriods = computed(() =>
      props.idleRule.disable_during_periods ? props.idleRule.disable_during_periods.join(', ') : '',
    );

    return {
      t,
      disableDuringPeriods,
      convertDateToString,
      convertDurationToString,
    };
  },
};
</script>
