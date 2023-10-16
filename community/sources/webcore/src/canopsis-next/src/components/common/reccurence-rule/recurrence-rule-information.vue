<template>
  <v-layout column>
    <c-information-block-row :label="$t('common.recurrence')">
      <span class="c-word-break-all">{{ recurrenceRuleString }}</span>
    </c-information-block-row>
    <c-information-block-row :label="$t('common.summary')">
      <span class="c-word-wrap-break">{{ recurrenceRuleText }}</span>
    </c-information-block-row>
  </v-layout>
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
