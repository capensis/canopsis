<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title
        span.headline {{ $t('modal.createDeclareTicket.title') }}
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
import AlarmGeneralTable from '@/components/tables/alarm-general.vue';
import EventEntityMixin from '@/mixins/event-entity';
import ModalItemMixin from '@/mixins/modal-item';
import { EVENT_TYPES } from '@/config';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    AlarmGeneralTable,
  },
  mixins: [ModalItemMixin, EventEntityMixin],
  methods: {
    async submit() {
      await this.createEvent(EVENT_TYPES.declareTicket, this.item, {
        output: 'declare ticket',
      });

      this.hideModal();
    },
  },
};
</script>
