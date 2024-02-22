<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ $t('common.confirmation') }}</span>
    </template>
    <template #text="">
      <v-alert type="info">
        {{ $t('modals.confirmAckWithTicket.infoMessage') }}
      </v-alert>
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="$modals.hide"
      >
        {{ $t('common.cancel') }}
      </v-btn>
      <v-btn
        :loading="submitting"
        :disabled="isDisabled || submittingWithTicket"
        class="primary"
        @click="submit"
      >
        {{ $t('common.continue') }}
      </v-btn>
      <v-btn
        :loading="submittingWithTicket"
        :disabled="isDisabledWithTicket || submitting"
        class="warning"
        @click="submitWithTicket"
      >
        {{ $t('modals.confirmAckWithTicket.continueAndAssociateTicket') }}
      </v-btn>
    </template>
  </modal-wrapper>
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
