<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createAssociateTicketEvent.title') }}</span>
      </template>
      <template #text="">
        <v-layout column>
          <alarm-general-table
            v-if="config.items"
            :items="config.items"
          />
          <v-divider />
          <v-checkbox
            v-if="isAllComponentAlarms"
            v-model="form.ticket_resources"
            :label="$t('alarm.associateTicketResources')"
            color="primary"
          />
          <associate-ticket-event-form
            v-model="form"
            class="mt-3"
          />
          <c-description-field
            v-model="form.comment"
            :label="$tc('common.comment')"
            name="comment"
          />
          <c-alert
            v-if="itemsWithoutAck.length"
            type="info"
          >
            {{ alertMessage }}
          </c-alert>
        </v-layout>
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
          {{ $t('common.saveChanges') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS } from '@/constants';

import { eventToAssociateTicketForm, formToAssociateTicketEvent } from '@/helpers/entities/associate-ticket/event/form';
import { isEntityComponentType } from '@/helpers/entities/entity/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';
import AssociateTicketEventForm from '@/components/other/declare-ticket/form/associate-ticket-event-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to associate a ticket to an alarm
 */
export default {
  name: MODALS.createAssociateTicketEvent,

  $_veeValidate: {
    validator: 'new',
  },
  components: { AssociateTicketEventForm, AlarmGeneralTable, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { associateTicket } = this.modal.config;

    return {
      form: eventToAssociateTicketForm(associateTicket),
    };
  },
  computed: {
    itemsWithoutAck() {
      return this.config.ignoreAck ? [] : this.config.items.filter(item => !item.v.ack);
    },

    alertMessage() {
      const { length: count } = this.itemsWithoutAck;

      return this.$tc('declareTicket.noAckItems', count, { count });
    },

    isAllComponentAlarms() {
      return this.config.items.every(({ entity }) => isEntityComponentType(entity.type));
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToAssociateTicketEvent(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
