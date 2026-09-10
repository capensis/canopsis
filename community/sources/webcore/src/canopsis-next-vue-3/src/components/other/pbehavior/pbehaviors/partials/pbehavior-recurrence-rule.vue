<template>
  <div>
    <recurrence-rule-information :rrule="pbehavior.rrule" />

    <pbehavior-recurrence-rule-periods :pbehavior="pbehavior" class="mt-2" />
  </div>
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
