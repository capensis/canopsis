<template lang="pug">
  v-card
    v-card-text
      c-information-block-row(:label="$t('common.description')") {{ eventFilter.description }}
      c-information-block-row(:label="$t('common.type')") {{ eventFilter.type }}
      c-information-block-row(:label="$t('common.priority')") {{ eventFilter.priority }}
      c-information-block-row(:label="$t('common.author')") {{ eventFilter.author.name }}
      c-information-block-row(:label="$t('common.created')") {{ eventFilter.created | date }}
      c-information-block-row(:label="$t('common.updated')") {{ eventFilter.updated | date }}
      c-information-block-row(v-if="eventFilter.start", :label="$t('common.start')") {{ eventFilter.start | date }}
      c-information-block-row(v-if="eventFilter.stop", :label="$t('common.stop')") {{ eventFilter.stop | date }}
      recurrence-rule-information(v-if="eventFilter.rrule", :rrule="eventFilter.rrule")
      pbehavior-exceptions-field(
        v-if="hasExdatesOrExceptions",
        :exdates="exdates",
        :exceptions="exceptions",
        disabled
      )
</template>

<script>
import { exceptionsToForm, exdatesToForm } from '@/helpers/forms/planning-pbehavior';

import RecurrenceRuleInformation from '@/components/common/reccurence-rule/recurrence-rule-information.vue';
import PbehaviorExceptionsField from '@/components/other/pbehavior/pbehaviors/fields/pbehavior-exceptions-field.vue';

export default {
  inject: ['$system'],
  components: { PbehaviorExceptionsField, RecurrenceRuleInformation },
  props: {
    eventFilter: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    hasExdatesOrExceptions() {
      return this.eventFilter.exdates?.length || this.eventFilter.exceptions?.length;
    },

    exdates() {
      return exdatesToForm(this.eventFilter.exdates, this.$system.timezone);
    },

    exceptions() {
      return exceptionsToForm(this.eventFilter.exceptions);
    },
  },
};
</script>
