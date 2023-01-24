<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createDeclareTicketEvent.title') }}
      template(#text="")
        declare-ticket-events-form(
          v-model="form",
          :alarms="items",
          :tickets-by-alarms="config.ticketsByAlarms",
          :alarms-by-tickets="config.alarmsByTickets"
        )
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

import { alarmsToDeclareTicketEventForm, formToDeclareTicketEvents } from '@/helpers/forms/declare-ticket-event';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { modalInnerItemsMixin } from '@/mixins/modal/inner-items';
import { eventActionsAlarmMixin } from '@/mixins/event-actions/alarm';
import { submittableMixinCreator } from '@/mixins/submittable';

import DeclareTicketEventsForm from '@/components/other/declare-ticket/form/declare-ticket-events-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to declare a ticket
 */
export default {
  name: MODALS.createDeclareTicketEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { DeclareTicketEventsForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    modalInnerItemsMixin,
    eventActionsAlarmMixin,
    submittableMixinCreator(),
  ],
  data() {
    const { alarmsByTickets } = this.modal.config;

    return {
      form: alarmsToDeclareTicketEventForm(alarmsByTickets),
    };
  },
  methods: {
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
