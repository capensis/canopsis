<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title
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
            :label="$t('modals.createAssociateTicket.fields.ticket')",
            :error-messages="errors.collect('ticket')",
            v-model="form.ticket",
            v-validate="'required'",
            data-vv-name="ticket"
            )
      v-card-actions
        v-btn(type="submit", :disabled="errors.any()", color="primary") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import AlarmGeneralTable from '@/components/other/alarm/alarm-general-list.vue';
import modalInnerItemsMixin from '@/mixins/modal/modal-inner-items';
import eventActionsMixin from '@/mixins/event-actions';
import { EVENT_ENTITY_TYPES, MODALS } from '@/constants';

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
  mixins: [modalInnerItemsMixin, eventActionsMixin],
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

        this.hideModal();
      }
    },
  },
};
</script>
