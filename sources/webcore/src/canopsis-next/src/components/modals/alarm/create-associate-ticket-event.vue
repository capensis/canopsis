<template lang="pug">
  v-form(data-test="createAssociateTicketModal", @submit.prevent="submit")
    v-card
      v-card-title.primary.white--text
        v-layout(justify-space-between, align-center)
          span.headline {{ $t('modals.createAssociateTicket.title') }}
      v-card-text
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table(:items="items")
          v-layout(row)
            v-divider.my-3
          v-layout(row)
            v-text-field(
              data-test="createAssociateTicketNumberOfTicket",
              :label="$t('modals.createAssociateTicket.fields.ticket')",
              :error-messages="errors.collect('ticket')",
              v-model="form.ticket",
              v-validate="'required'",
              data-vv-name="ticket"
            )
      v-divider
      v-layout.py-1(justify-end)
        v-btn(
          data-test="createAssociateTicketCancelButton",
          @click="$modals.hide",
          depressed,
          flat
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          data-test="createAssociateTicketSubmitButton",
          type="submit",
          :disabled="errors.any()"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS, EVENT_ENTITY_TYPES } from '@/constants';

import AlarmGeneralTable from '@/components/other/alarm/alarm-general-list.vue';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';

/**
 * Modal to associate a ticket to an alarm
 */
export default {
  name: MODALS.createAssociateTicketEvent,

  $_veeValidate: {
    validator: 'new',
  },
  components: {
    AlarmGeneralTable,
  },
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin],
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
