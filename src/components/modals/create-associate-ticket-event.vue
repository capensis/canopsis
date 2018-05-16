<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title
        span.headline {{ $t('modals.createAssociateTicket.title') }}
      v-card-text
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table(:item="item")
          v-layout(row)
            v-divider.my-3
          v-layout(row)
            v-text-field(
            :label="$t('modals.createAssociateTicket.ticket')",
            :error-messages="errors.collect('ticket')",
            v-model="form.ticket",
            v-validate="'required'",
            data-vv-name="ticket"
            )
      v-card-actions
        v-btn(type="submit", :disabled="errors.any()", color="primary") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import AlarmGeneralTable from '@/components/tables/alarm-general.vue';
import EventEntityMixin from '@/mixins/event-entity';
import ModalItemMixin from '@/mixins/modal-item';
import { EVENT_TYPES } from '@/config';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    AlarmGeneralTable,
  },
  mixins: [ModalItemMixin, EventEntityMixin],
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
        await this.createEvent(EVENT_TYPES.associateTicket, this.item, this.form);

        this.hideModal();
      }
    },
  },
};
</script>
