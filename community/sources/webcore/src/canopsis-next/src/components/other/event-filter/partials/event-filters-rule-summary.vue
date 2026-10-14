<template>
  <c-layout column>
    <c-information-block-row :label="$t('common.description')" :width="labelWidth">
      {{ eventFilter.description }}
    </c-information-block-row>
    <c-information-block-row :label="$t('common.type')" :width="labelWidth">
      {{ $t(`eventFilter.types.${eventFilter.type}`) }}
    </c-information-block-row>
    <c-information-block-row :label="$t('common.priority')" :width="labelWidth">
      {{ eventFilter.priority || '-' }}
    </c-information-block-row>
    <c-information-block-row :label="$t('common.author')" :width="labelWidth">
      {{ eventFilter.author?.display_name }}
    </c-information-block-row>
    <c-information-block-row :label="$t('common.created')" :width="labelWidth">
      {{ dateToString(eventFilter.created) }}
    </c-information-block-row>
    <c-information-block-row :label="$t('common.updated')" :width="labelWidth">
      {{ dateToString(eventFilter.updated) }}
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
    <div v-if="hasExecutionDate" class="text-subtitle-1 font-weight-bold mt-4">
      {{ $t('eventFilter.executionDate') }}
    </div>
    <c-information-block-row
      v-if="eventFilter.start"
      :label="$t('common.start')"
      :width="labelWidth"
    >
      {{ dateToString(eventFilter.start) }}
    </c-information-block-row>
    <c-information-block-row v-if="eventFilter.stop" :label="$t('common.stop')" :width="labelWidth">
      {{ dateToString(eventFilter.stop) }}
    </c-information-block-row>
    <recurrence-rule-information v-if="eventFilter.rrule" :rrule="eventFilter.rrule" />
    <c-layout v-if="exceptions.length" class="mt-3" column>
      <div class="text-subtitle-1 font-weight-bold">
        {{ $t('pbehavior.exceptions.title') }}
      </div>
      <pbehavior-exceptions-list :exceptions="exceptions" disabled />
    </c-layout>
    <c-layout v-if="exdates.length" class="mt-3" column>
      <div class="text-subtitle-1 font-weight-bold">
        {{ $t('pbehavior.exdates.title') }}
      </div>
      <pbehavior-exceptions-field :exdates="exdates" disabled />
    </c-layout>
  </c-layout>
</template>

<script>
import { convertDateToString } from '@/helpers/date/date';
import { exceptionsToForm, exdatesToForm } from '@/helpers/entities/pbehavior/form';

import RecurrenceRuleInformation from '@/components/common/reccurence-rule/recurrence-rule-information.vue';
import PbehaviorExceptionsField from '@/components/other/pbehavior/exceptions/fields/pbehavior-exceptions-field.vue';
import PbehaviorExceptionsList from '@/components/other/pbehavior/pbehaviors/partials/pbehavior-exceptions-list.vue';

export default {
  components: {
    RecurrenceRuleInformation,
    PbehaviorExceptionsField,
    PbehaviorExceptionsList,
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
    hasExecutionDate() {
      return this.eventFilter.start ||
        this.eventFilter.stop ||
        this.eventFilter.rrule ||
        this.hasExdatesOrExceptions;
    },

    hasExdatesOrExceptions() {
      return this.eventFilter.exdates?.length || this.eventFilter.exceptions?.length;
    },

    exdates() {
      return exdatesToForm(this.eventFilter.exdates);
    },

    exceptions() {
      return exceptionsToForm(this.eventFilter.exceptions);
    },
  },
  methods: {
    dateToString(value) {
      return convertDateToString(value);
    },
  },
};
</script>
