<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ config.pbehavior.name }}
    template(#text="")
      pbehavior-recurrence-rule-calendar(:pbehavior="config.pbehavior")
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import PbehaviorRecurrenceRuleCalendar from '@/components/other/pbehavior/calendar/pbehavior-recurrence-rule-calendar.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.pbehaviorRecurrenceRule,
  components: { PbehaviorRecurrenceRuleCalendar, ModalWrapper },
  mixins: [modalInnerMixin],
};
</script>
