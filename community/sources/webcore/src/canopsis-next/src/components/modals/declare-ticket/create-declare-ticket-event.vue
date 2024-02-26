<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createDeclareTicketEvent.title') }}</span>
      </template>
      <template #text="">
        <c-enabled-field
          v-if="config.items.length > 1"
          v-model="singleMode"
          :label="$t('declareTicket.oneTicketForAlarms')"
        />
        <declare-ticket-events-form
          v-model="form"
          :alarms="config.items"
          :tickets-by-alarms="config.ticketsByAlarms"
          :alarms-by-tickets="config.alarmsByTickets"
          :hide-ticket-resource="!isAllComponentAlarms"
          :single-mode="singleMode"
        />
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
          :disabled="isDisabled"
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import {
  alarmsToDeclareTicketEventForm,
  formToDeclareTicketEvents,
} from '@/helpers/entities/declare-ticket/event/form';
import { isEntityComponentType } from '@/helpers/entities/entity/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import DeclareTicketEventsForm from '@/components/other/declare-ticket/form/declare-ticket-events-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to declare a ticket
 */
export default {
  name: MODALS.createDeclareTicketEvent,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { DeclareTicketEventsForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { alarmsByTickets } = this.modal.config;

    return {
      singleMode: false,
      form: alarmsToDeclareTicketEventForm(alarmsByTickets),
    };
  },
  computed: {
    isAllComponentAlarms() {
      return this.config.items.every(({ entity }) => isEntityComponentType(entity.type));
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToDeclareTicketEvents(this.form, this.singleMode), this.singleMode);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
