<template lang="pug">
  v-form(data-test="createAssociateTicketModal", @submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createAssociateTicket.title') }}
      template(slot="text")
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
              name="ticket",
              data-test="createAssociateTicketNumberOfTicket"
            )
      template(slot="actions")
        v-btn(
          data-test="createAssociateTicketCancelButton",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          data-test="createAssociateTicketSubmitButton",
          type="submit"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS, EVENT_ENTITY_TYPES } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';
import submittableMixin from '@/mixins/submittable';


import AlarmGeneralTable from '@/components/other/alarm/alarm-general-list.vue';

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
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin, submittableMixin()],
  data() {
    return {
      form: {
        ticket: '',
        output: 'Associated ticket number',
      },
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.createEvent(EVENT_ENTITY_TYPES.assocTicket, this.items, this.form);

        this.$modals.hide();
      }
    },
  },
};
</script>
