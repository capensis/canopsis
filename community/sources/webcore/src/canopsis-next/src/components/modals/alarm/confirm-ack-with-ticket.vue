<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('common.confirmation') }}
    template(#text="")
      v-alert(:value="true", type="info") {{ $t('modals.confirmAckWithTicket.infoMessage') }}
    template(#actions="")
      v-btn(
        depressed,
        flat,
        @click="$modals.hide"
      ) {{ $t('common.cancel') }}
      v-btn.primary(
        :loading="submitting",
        :disabled="isDisabled || submittingWithTicket",
        @click="submit"
      ) {{ $t('common.continue') }}
      v-btn.warning(
        :loading="submittingWithTicket",
        :disabled="isDisabledWithTicket || submitting",
        @click="submitWithTicket"
      ) {{ $t('modals.confirmAckWithTicket.continueAndAssociateTicket') }}
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Ack with ticket confirmation modal
 */
export default {
  name: MODALS.confirmAckWithTicket,
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    submittableMixinCreator({
      method: 'submitWithTicket',
      property: 'submittingWithTicket',
      computedProperty: 'isDisabledWithTicket',
    }),
  ],
  data() {
    return {
      submitting: false,
    };
  },
  methods: {
    async submit() {
      if (this.config.continueAction) {
        await this.config.continueAction();
      }

      this.$modals.hide();
    },

    async submitWithTicket() {
      if (this.config.continueWithTicketAction) {
        await this.config.continueWithTicketAction();
      }

      this.$modals.hide();
    },
  },
};
</script>
