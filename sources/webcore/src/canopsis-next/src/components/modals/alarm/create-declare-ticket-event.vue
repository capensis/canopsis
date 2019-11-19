<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('modals.createDeclareTicket.title') }}
    template(slot="text")
      v-container
        v-layout(row)
          v-flex.text-xs-center
            alarm-general-table(:items="items")
        v-layout(row)
          v-divider.my-3
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn.primary(:disabled="submitting", @click="submit") {{ $t('common.actions.reportIncident') }}
</template>

<script>
import { MODALS, EVENT_ENTITY_TYPES } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';

import AlarmGeneralTable from '@/components/other/alarm/alarm-general-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to declare a ticket
 */
export default {
  name: MODALS.createDeclareTicketEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AlarmGeneralTable, ModalWrapper },
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin],
  data() {
    return {
      submitting: false,
    };
  },
  methods: {
    async submit() {
      try {
        this.submitting = true;

        await this.createEvent(EVENT_ENTITY_TYPES.declareTicket, this.items, {
          output: 'declare ticket',
        });

        this.$modals.hide();
      } catch (err) {
        this.$popups.error({ text: err.description || this.$t('error.default') });
      } finally {
        this.submitting = false;
      }
    },
  },
};
</script>
