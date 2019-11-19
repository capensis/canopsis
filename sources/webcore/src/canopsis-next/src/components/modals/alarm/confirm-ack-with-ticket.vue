<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('common.confirmation') }}
    template(slot="text")
      v-alert(:value="true", type="info") {{ $t('modals.confirmAckWithTicket.infoMessage') }}
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn.primary(
        :loading="submitting",
        :disabled="submitting",
        @click="submit"
      ) {{ $t('common.continue') }}
      v-btn.warning(
        :loading="submitting",
        :disabled="submitting",
        @click="submitWithTicket"
      ) {{ $t('modals.confirmAckWithTicket.continueAndAssociateTicket') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Ack with ticket confirmation modal
 */
export default {
  name: MODALS.confirmAckWithTicket,
  components: { ModalWrapper },
  mixins: [modalInnerMixin],
  data() {
    return {
      submitting: false,
    };
  },
  methods: {
    async submit() {
      try {
        this.submitting = true;

        if (this.config.continueAction) {
          await this.config.continueAction();
        }

        this.$modals.hide();
      } catch (err) {
        this.$popups.error({ text: err.description || this.$t('error.default') });
      } finally {
        this.submitting = false;
      }
    },

    async submitWithTicket() {
      try {
        this.submitting = true;

        if (this.config.continueWithTicketAction) {
          await this.config.continueWithTicketAction();
        }

        this.$modals.hide();
      } catch (err) {
        this.$popups.error({ text: err.description || this.$t('error.default') });
      } finally {
        this.submitting = false;
      }
    },
  },
};
</script>
