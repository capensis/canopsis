<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createAssociateTicket.title') }}
      template(#text="")
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table(:items="items")
          v-layout(row)
            v-divider.my-3
          v-layout(row)
            v-text-field(
              v-model="form.ticket",
              v-validate="'required'",
              :label="$t('modals.createAssociateTicket.fields.ticket')",
              :error-messages="errors.collect('ticket')",
              name="ticket"
            )
          v-alert(:value="itemsWithoutAck.length", type="info")
            span {{ alertMessage }}
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
        ) {{ $t('common.saveChanges') }}
</template>

<script>
import { MODALS, EVENT_ENTITY_TYPES } from '@/constants';

import { prepareEventsByAlarms } from '@/helpers/forms/event';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { modalInnerItemsMixin } from '@/mixins/modal/inner-items';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import { eventActionsAlarmMixin } from '@/mixins/event-actions/alarm';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to associate a ticket to an alarm
 */
export default {
  name: MODALS.createAssociateTicketEvent,

  $_veeValidate: {
    validator: 'new',
  },
  components: { AlarmGeneralTable, ModalWrapper },
  mixins: [
    modalInnerMixin,
    modalInnerItemsMixin,
    eventActionsAlarmMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: {
        ticket: '',
        output: 'Associated ticket number',
      },
    };
  },
  computed: {
    itemsWithoutAck() {
      return this.items.filter(item => !item.v.ack);
    },

    alertMessage() {
      const { length: count } = this.itemsWithoutAck;

      return this.$tc('modals.createAssociateTicket.alerts.noAckItems', count, { count });
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.itemsWithoutAck.length) {
          const { fastAckOutput } = this.config;

          const eventData = prepareEventsByAlarms(EVENT_ENTITY_TYPES.ack, this.itemsWithoutAck, {
            output: fastAckOutput && fastAckOutput.enabled ? fastAckOutput.value : '',
            ticket: this.form.ticket,
          });

          await this.createEventAction({ data: eventData });
        }

        await this.createEvent(EVENT_ENTITY_TYPES.assocTicket, this.items, this.form);

        this.$modals.hide();
      }
    },
  },
};
</script>
