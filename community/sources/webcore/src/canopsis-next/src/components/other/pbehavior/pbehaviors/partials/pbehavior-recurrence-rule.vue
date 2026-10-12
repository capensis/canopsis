<template>
  <v-layout>
    <v-flex
      offset-xs3
      xs6
    >
      <recurrence-rule-information :rrule="pbehavior.rrule" />
      <pbehavior-recurrence-rule-periods
        :pbehavior="pbehavior"
        class="mt-2"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import RecurrenceRuleInformation from '@/components/common/reccurence-rule/recurrence-rule-information.vue';

import PbehaviorRecurrenceRulePeriods from './pbehavior-recurrence-rule-periods.vue';

export default {
  components: { RecurrenceRuleInformation, PbehaviorRecurrenceRulePeriods },
  props: {
    pbehavior: {
      type: Object,
      required: true,
    },
  },
};
</script>
