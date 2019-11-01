<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('common.confirmation') }}
    v-card-text
      v-alert(:value="true", type="info") {{ $t('modals.confirmAckWithTicket.infoMessage') }}
      v-divider
      v-layout.mt-2.mb-1(wrap, justify-end)
        v-btn(@click="$modals.hide", flat) {{ $t('common.cancel') }}
        v-btn(
          @click.prevent="submit",
          :loading="submitting",
          :disabled="submitting",
          color="primary"
        ) {{ $t('common.continue') }}
        v-btn(
          @click.prevent="submitWithTicket",
          :loading="submitting",
          :disabled="submitting",
          color="warning"
        ) {{ $t('modals.confirmAckWithTicket.continueAndAssociateTicket') }}
</template>

<script>
import modalInnerMixin from '@/mixins/modal/inner';
import { MODALS } from '@/constants';

/**
 * Ack with ticket confirmation modal
 */
export default {
  name: MODALS.confirmAckWithTicket,
  mixins: [modalInnerMixin],
  data() {
    return {
      submitting: false,
    };
  },
  methods: {
    async submit() {
      this.submitting = true;
      if (this.config.continueAction) {
        await this.config.continueAction();
      }
      this.$modals.hide();
      this.submitting = false;
    },

    async submitWithTicket() {
      this.submitting = true;
      if (this.config.continueWithTicketAction) {
        await this.config.continueWithTicketAction();
      }
      this.$modals.hide();
      this.submitting = false;
    },
  },
};
</script>
