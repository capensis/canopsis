<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title
        span.headline {{ $t('modals.createDeclareTicket.title') }}
      v-card-text
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table(:item="item")
          v-layout(row)
            v-divider.my-3
      v-card-actions
        v-btn(type="submit", color="warning") {{ $t('common.actions.reportIncident') }}
</template>

<script>
import AlarmGeneralTable from '@/components/tables/alarm/general.vue';
import ModalInnerItemMixin from '@/mixins/modal/modal-inner-item';
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
  mixins: [ModalInnerItemMixin, EventActionsMixin],
  methods: {
    async submit() {
      await this.createEvent(EVENT_ENTITY_TYPES.declareTicket, this.item, {
        output: 'declare ticket',
      });

      this.hideModal();
    },
  },
};
</script>
