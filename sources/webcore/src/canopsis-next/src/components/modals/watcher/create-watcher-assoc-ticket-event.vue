<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createDeclareTicket.title') }}
        v-btn(icon, dark, @click="hideModal")
          v-icon close
    v-card-text
      v-text-field(
      :label="$t('modals.createAssociateTicket.fields.ticket')",
      v-model="ticket",
      name="ticket",
      v-validate="'required'",
      :error-messages="errors.collect('ticket')"
      )
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(type="submit", @click="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

export default {
  name: MODALS.createWatcherAssocTicketEvent,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      ticket: '',
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.ticket);
        }

        this.hideModal();
      }
    },
  },
};
</script>
