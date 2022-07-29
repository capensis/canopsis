<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createAckEvent.title') }}
      template(#text="")
        v-container
          v-layout(row, align-center)
            v-flex.text-xs-center
              alarm-general-table(:items="items")
          v-layout(row)
            v-divider.my-3
          ack-event-form(v-model="form", :is-note-required="isNoteRequired")
      template(#actions="")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled || submittingWithTicket",
          type="submit"
        ) {{ $t('common.acknowledge') }}
        v-btn.warning(
          :loading="submittingWithTicket",
          :disabled="isDisabledWithTicket || submitting",
          @click="submitWithTicket"
        ) {{ submitWithTicketBtnLabel }}
</template>

<script>
import { MODALS, EVENT_ENTITY_TYPES } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { modalInnerItemsMixin } from '@/mixins/modal/inner-items';
import { eventActionsAlarmMixin } from '@/mixins/event-actions/alarm';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import { prepareEventsByAlarms } from '@/helpers/forms/event';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';
import AckEventForm from '@/components/widgets/alarm/forms/ack-event-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create an ack event
 */
export default {
  name: MODALS.createAckEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AlarmGeneralTable, AckEventForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    modalInnerItemsMixin,
    eventActionsAlarmMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
    submittableMixinCreator({
      method: 'submitWithTicket',
      property: 'submittingWithTicket',
      computedProperty: 'isDisabledWithTicket',
    }),
  ],
  data() {
    return {
      form: {
        ticket: '',
        output: '',
        ack_resources: false,
      },
    };
  },
  computed: {
    isNoteRequired() {
      return this.config?.isNoteRequired;
    },

    submitWithTicketBtnLabel() {
      return this.form.ticket
        ? this.$t('common.acknowledgeAndAssociateTicket')
        : this.$t('common.acknowledgeAndDeclareTicket');
    },
  },
  methods: {
    createAckEvent() {
      const ackEventData = prepareEventsByAlarms(EVENT_ENTITY_TYPES.ack, this.items, this.form);

      return this.createEventAction({ data: ackEventData });
    },

    createDeclareTicketEvent() {
      const declareTicketEventData = prepareEventsByAlarms(
        EVENT_ENTITY_TYPES.declareTicket,
        this.items,
        { output: this.form.output },
      );

      return this.createEventAction({ data: declareTicketEventData });
    },

    createAssocTicketEvent() {
      const assocTicketEventData = prepareEventsByAlarms(
        EVENT_ENTITY_TYPES.assocTicket,
        this.items,
        { ticket: this.form.ticket, output: this.form.output },
      );

      return this.createEventAction({ data: assocTicketEventData });
    },

    async createAckEventAndCloseModal() {
      await this.createAckEvent();

      if (this.config && this.config.afterSubmit) {
        await this.config.afterSubmit();
      }

      this.$modals.hide();
    },

    async submitWithTicket() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.form.ticket) {
          await this.createAssocTicketEvent();
        } else {
          await this.createDeclareTicketEvent();
        }

        await this.createAckEventAndCloseModal();
      }
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.form.ticket) {
          this.$modals.show({
            name: MODALS.confirmAckWithTicket,
            config: {
              continueAction: this.createAckEventAndCloseModal,
              continueWithTicketAction: this.submitWithTicket,
            },
          });
        } else {
          await this.createAckEventAndCloseModal();
        }
      }
    },
  },
};
</script>
