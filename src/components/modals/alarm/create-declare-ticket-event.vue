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
import AlarmGeneralTable from '@/components/other/alarm/alarm-general-list.vue';
import modalInnerItemsMixin from '@/mixins/modal/modal-inner-items';
import eventActionsMixin from '@/mixins/event-actions';
import { MODALS } from '@/constants';

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
  mixins: [modalInnerItemsMixin, eventActionsMixin],
  methods: {
    async submit() {
      await this.createEvent(this.$constants.EVENT_ENTITY_TYPES.declareTicket, this.items, {
        output: 'declare ticket',
      });

      this.hideModal();
    },
  },
};
</script>
