<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createAssociateTicketEvent.title') }}
      template(#text="")
        v-layout(column)
          alarm-general-table(:items="items")
          v-divider
          associate-ticket-event-form.mt-3(v-model="form")
          c-alert(v-if="itemsWithoutAck.length", type="info") {{ alertMessage }}
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
import { MODALS } from '@/constants';

import { eventToAssociateTicketForm, formToAssociateTicketEvent } from '@/helpers/forms/associate-ticket-event';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { modalInnerItemsMixin } from '@/mixins/modal/inner-items';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';
import AssociateTicketEventForm from '@/components/other/declare-ticket/form/associate-ticket-event-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to associate a ticket to an alarm
 */
export default {
  name: MODALS.createAssociateTicketEvent,

  $_veeValidate: {
    validator: 'new',
  },
  components: { AssociateTicketEventForm, AlarmGeneralTable, ModalWrapper },
  mixins: [
    modalInnerMixin,
    modalInnerItemsMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { associateTicket } = this.modal.config;

    return {
      form: eventToAssociateTicketForm(associateTicket),
    };
  },
  computed: {
    itemsWithoutAck() {
      return this.items.filter(item => !item.v.ack);
    },

    alertMessage() {
      const { length: count } = this.itemsWithoutAck;

      return this.$tc('declareTicket.noAckItems', count, { count });
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToAssociateTicketEvent(this.form));
        }

        if (this.config?.afterSubmit) {
          this.config.afterSubmit();
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
