<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createAckEvent.title') }}</span>
      </template>
      <template #text="">
        <v-layout column>
          <alarm-general-table
            v-if="config.items"
            :items="config.items"
          />
          <v-divider />
          <ack-event-form
            v-model="form"
            :is-note-required="isNoteRequired"
            :hide-ack-resources="!isAllComponentAlarms"
          />
          <v-radio-group
            v-model="actionType"
            :label="$t('alarm.actionsRequired')"
            name="actionType"
            hide-details
            column
          >
            <v-radio
              v-for="type in actionTypes"
              :key="type.value"
              :value="type.value"
              :label="type.label"
              color="primary"
            />
          </v-radio-group>
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
          class="primary"
          :disabled="isDisabled"
          :loading="submitting"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, ACK_MODAL_ACTIONS_TYPES } from '@/constants';

import { isEntityComponentType } from '@/helpers/entities/entity/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';
import AckEventForm from '@/components/widgets/alarm/forms/ack-event-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create an ack event
 */
export default {
  name: MODALS.createAckEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AlarmGeneralTable, AckEventForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      actionType: ACK_MODAL_ACTIONS_TYPES.ack,
      form: {
        comment: '',
        ack_resources: false,
      },
    };
  },
  computed: {
    isAlarmsHasDeclareTicketRules() {
      return this.config.items.some(
        ({ assigned_declare_ticket_rules: assignedDeclareTicketRules }) => assignedDeclareTicketRules?.length,
      );
    },

    isAllComponentAlarms() {
      return this.config.items.every(({ entity }) => isEntityComponentType(entity.type));
    },

    actionTypes() {
      const types = [
        {
          label: this.$t('alarm.acknowledge'),
          value: ACK_MODAL_ACTIONS_TYPES.ack,
        },
        {
          label: this.$t('alarm.acknowledgeAndAssociateTicket'),
          value: ACK_MODAL_ACTIONS_TYPES.ackAndAssociateTicket,
        },
      ];

      if (this.isAlarmsHasDeclareTicketRules) {
        types.push({
          label: this.$t('alarm.acknowledgeAndDeclareTicket'),
          value: ACK_MODAL_ACTIONS_TYPES.ackAndDeclareTicket,
        });
      }

      return types;
    },

    isNoteRequired() {
      return this.config?.isNoteRequired;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          const needDeclareTicket = this.actionType === ACK_MODAL_ACTIONS_TYPES.ackAndDeclareTicket;
          const needAssociateTicket = this.actionType === ACK_MODAL_ACTIONS_TYPES.ackAndAssociateTicket;

          await this.config.action(this.form, {
            needDeclareTicket,
            needAssociateTicket,
          });
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
