<template>
  <v-layout column>
    <c-information-block-row
      :label="$t('common.description')"
      :width="labelWidth"
    >
      {{ eventFilter.description }}
    </c-information-block-row>
    <c-information-block-row
      :label="$t('common.type')"
      :width="labelWidth"
    >
      {{ $t(`eventFilter.types.${eventFilter.type}`) }}
    </c-information-block-row>
    <c-information-block-row
      :label="$t('common.priority')"
      :width="labelWidth"
    >
      {{ eventFilter.priority || '-' }}
    </c-information-block-row>
    <c-information-block-row
      :label="$t('common.author')"
      :width="labelWidth"
    >
      {{ eventFilter.author?.display_name }}
    </c-information-block-row>
    <c-information-block-row
      :label="$t('common.created')"
      :width="labelWidth"
    >
      {{ eventFilter.created | date }}
    </c-information-block-row>
    <c-information-block-row
      :label="$t('common.updated')"
      :width="labelWidth"
    >
      {{ eventFilter.updated | date }}
    </c-information-block-row>
    <c-information-block-row
      v-if="eventFilter.events_count"
      :label="$t('eventFilter.eventsFilteredSinceLastUpdate')"
      :width="labelWidth"
    >
      {{ eventFilter.events_count }}
    </c-information-block-row>
    <c-information-block-row
      v-if="eventFilter.failures_count"
      :label="$t('eventFilter.errorsSinceLastUpdate')"
      :width="labelWidth"
    >
      {{ eventFilter.failures_count }}
    </c-information-block-row>
    <c-information-block-row
      v-if="eventFilter.start"
      :label="$t('common.start')"
      :width="labelWidth"
    >
      {{ eventFilter.start | date }}
    </c-information-block-row>
    <c-information-block-row
      v-if="eventFilter.stop"
      :label="$t('common.stop')"
      :width="labelWidth"
    >
      {{ eventFilter.stop | date }}
    </c-information-block-row>
    <recurrence-rule-information
      v-if="eventFilter.rrule"
      :rrule="eventFilter.rrule"
    />
    <pbehavior-recurrence-rule-exceptions-field
      v-if="hasExdatesOrExceptions"
      :exdates="exdates"
      :exceptions="exceptions"
      disabled
    />
  </v-layout>
</template>

<script>
import { exceptionsToForm, exdatesToForm } from '@/helpers/entities/pbehavior/form';

import RecurrenceRuleInformation from '@/components/common/reccurence-rule/recurrence-rule-information.vue';
import PbehaviorRecurrenceRuleExceptionsField from '@/components/other/pbehavior/exceptions/fields/pbehavior-recurrence-rule-exceptions-field.vue';

export default {
  inject: ['$system'],
  components: {
    RecurrenceRuleInformation,
    PbehaviorRecurrenceRuleExceptionsField,
  },
  props: {
    eventFilter: {
      type: Object,
      default: () => ({}),
    },
    labelWidth: {
      type: [Number, String],
      default: 220,
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
