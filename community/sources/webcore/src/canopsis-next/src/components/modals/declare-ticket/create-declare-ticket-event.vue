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
          :alarms-by-tickets="config.alarmsByTickets",
          :hide-ticket-resource="!isAllComponentAlarms"
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
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { alarmsToDeclareTicketEventForm, formToDeclareTicketEvents } from '@/helpers/forms/declare-ticket-event';
import { isEntityComponentType } from '@/helpers/entities/entity';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { modalInnerItemsMixin } from '@/mixins/modal/inner-items';
import { eventActionsAlarmMixin } from '@/mixins/event-actions/alarm';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import DeclareTicketEventsForm from '@/components/other/declare-ticket/form/declare-ticket-events-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to declare a ticket
 */
export default {
  name: MODALS.createDeclareTicketEvent,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { DeclareTicketEventsForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    modalInnerItemsMixin,
    eventActionsAlarmMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { alarmsByTickets } = this.modal.config;

    return {
      form: alarmsToDeclareTicketEventForm(alarmsByTickets),
    };
  },
  computed: {
    isAllComponentAlarms() {
      return this.items.every(({ entity }) => isEntityComponentType(entity.type));
    },
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
