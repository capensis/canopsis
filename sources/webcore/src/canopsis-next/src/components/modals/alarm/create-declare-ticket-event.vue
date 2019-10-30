<template lang="pug">
  v-form(date-test="createDeclareTicketEventModal", @submit.prevent="submit")
    v-card
      v-card-title.primary.white--text
        v-layout(justify-space-between, align-center)
          span.headline {{ $t('modals.createDeclareTicket.title') }}
      v-card-text
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table(:items="items")
          v-layout(row)
            v-divider.my-3
      v-divider
      v-layout.py-1(justify-end)
        v-btn(
          data-test="declareTicketEventCancelButton",
          @click="hideModal",
          depressed,
          flat
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          data-test="declareTicketEventSubmitButton",
          type="submit"
        ) {{ $t('common.actions.reportIncident') }}
</template>

<script>
import { MODALS, EVENT_ENTITY_TYPES } from '@/constants';

import AlarmGeneralTable from '@/components/other/alarm/alarm-general-list.vue';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';

/**
 * Modal to declare a ticket
 */
export default {
  name: MODALS.createDeclareTicketEvent,

  $_veeValidate: {
    validator: 'new',
  },
  components: {
    AlarmGeneralTable,
  },
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin],
  methods: {
    async submit() {
      await this.createEvent(EVENT_ENTITY_TYPES.declareTicket, this.items, {
        output: 'declare ticket',
      });

      this.hideModal();
    },
  },
};
</script>
