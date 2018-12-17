<template lang="pug">
  v-form(@submit.prevent="submit")
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
            :label="$t('modals.createAssociateTicket.fields.ticket')",
            :error-messages="errors.collect('ticket')",
            v-model="form.ticket",
            v-validate="'required'",
            data-vv-name="ticket"
            )
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(type="submit", :disabled="errors.any()") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import AlarmGeneralTable from '@/components/other/alarm/alarm-general-list.vue';
import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsMixin from '@/mixins/event-actions';
import { MODALS } from '@/constants';

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
        await this.createEvent(this.$constants.EVENT_ENTITY_TYPES.assocTicket, this.items, this.form);

        this.hideModal();
      }
    },
  },
};
</script>
