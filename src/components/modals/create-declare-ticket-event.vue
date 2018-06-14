<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title
        span.headline {{ $t('modals.createDeclareTicket.title') }}
      v-card-text
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table(:items="items")
          v-layout(row)
            v-divider.my-3
      v-card-actions
        v-btn(type="submit", color="warning") {{ $t('common.actions.reportIncident') }}
</template>

<script>
import AlarmGeneralTable from '@/components/other/alarm-list/alarm-general-list.vue';
import ModalInnerItemsMixin from '@/mixins/modal/modal-inner-items';
import EventActionsMixin from '@/mixins/event-actions';
import { EVENT_ENTITY_TYPES, MODALS } from '@/constants';

export default {
  name: MODALS.createDeclareTicketEvent,

  $_veeValidate: {
    validator: 'new',
  },
  components: {
    AlarmGeneralTable,
  },
  mixins: [ModalInnerItemsMixin, EventActionsMixin],
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
