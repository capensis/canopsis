<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createDeclareTicketEvent.title') }}
      template(#text="")
        c-progress-overlay(:pending="pending")
        declare-ticket-event-form(v-model="form", :alarms="items", :tickets-by-alarms="ticketsByAlarms")
      template(#actions="")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { declareTicketEventToForm, formToDeclareTicketEvents } from '@/helpers/forms/declare-ticket-event';
import { mapIds } from '@/helpers/entities';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { modalInnerItemsMixin } from '@/mixins/modal/inner-items';
import { eventActionsAlarmMixin } from '@/mixins/event-actions/alarm';
import { entitiesDeclareTicketRuleMixin } from '@/mixins/entities/declare-ticket-rule';
import { submittableMixinCreator } from '@/mixins/submittable';

import DeclareTicketEventForm from '@/components/other/declare-ticket/form/declare-ticket-event-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to declare a ticket
 */
export default {
  name: MODALS.createDeclareTicketEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { DeclareTicketEventForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    modalInnerItemsMixin,
    eventActionsAlarmMixin,
    entitiesDeclareTicketRuleMixin,
    submittableMixinCreator(),
  ],
  data() {
    const { declareTicketEvent, items } = this.modal.config;

    return {
      form: declareTicketEventToForm(declareTicketEvent, items),
      pending: true,
      ticketsByAlarms: {},
    };
  },
  mounted() {
    this.fetchAlarmsTickets();
  },
  methods: {
    async fetchAlarmsTickets() {
      this.pending = true;

      const { by_alarms: ticketsByAlarms } = await this.fetchAssignedDeclareTicketsWithoutStore({
        params: {
          alarms: mapIds(this.config.items),
        },
      });

      this.ticketsByAlarms = ticketsByAlarms;

      this.pending = false;
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToDeclareTicketEvents(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
