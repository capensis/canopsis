<template>
  <v-layout column>
    <idle-rules-summary-row
      :label="$t('common.id')"
      :value="idleRule._id"
    />
    <idle-rules-summary-row
      :label="$t('common.name')"
      :value="idleRule.name"
    />
    <idle-rules-summary-row
      :label="$t('common.description')"
      :value="idleRule.description"
    />
    <idle-rules-summary-row
      :label="$t('common.type')"
      :value="$t(`idleRules.types.${idleRule.type}`)"
    />
    <idle-rules-summary-row
      :label="$t('common.created')"
      :value="idleRule.created | date"
    />
    <idle-rules-summary-row
      :label="$t('common.updated')"
      :value="idleRule.updated | date"
    />
    <idle-rules-summary-row
      :label="$t('idleRules.timeAwaiting')"
      :value="idleRule.duration | duration"
    />
    <idle-rules-summary-row
      :label="$t('common.priority')"
      :value="idleRule.priority"
    />
    <idle-rules-summary-row
      v-if="disableDuringPeriods"
      :label="$t('common.disableDuringPeriods')"
      :value="disableDuringPeriods"
    />
  </v-layout>
</template>

<script>
import IdleRulesSummaryRow from './idle-rules-summary-row.vue';

export default {
  components: { IdleRulesSummaryRow },
  props: {
    idleRule: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    disableDuringPeriods() {
      return this.idleRule.disable_during_periods ? this.idleRule.disable_during_periods.join(', ') : '';
    },
  },
};
</script>
