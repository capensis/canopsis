<template lang="pug">
  v-layout(column)
    c-information-block-row(:label="$t('common.recurrence')")
      span.c-word-break-all {{ recurrenceRuleString }}
    c-information-block-row(:label="$t('common.summary')")
      span.c-word-wrap-break {{ recurrenceRuleText }}
</template>

<script>
import { rrulestr } from 'rrule';

export default {
  props: {
    rrule: {
      type: String,
      required: true,
    },
  },
  computed: {
    recurrenceRule() {
      return rrulestr(this.rrule);
    },

    recurrenceRuleString() {
      return this.recurrenceRule.toString();
    },

    recurrenceRuleText() {
      return this.recurrenceRule.toText();
    },
  },
};
</script>
