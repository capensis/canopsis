<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createAckEvent.title') }}
      template(#text="")
        v-layout(column)
          alarm-general-table(v-if="config.items", :items="config.items")
          v-divider
          ack-event-form(v-model="form", :is-note-required="isNoteRequired")
          v-radio-group(
            v-model="actionType",
            :label="$t('alarm.actionsRequired')",
            name="actionType",
            hide-details,
            column
          )
            v-radio(
              v-for="type in actionTypes",
              :key="type.value",
              :value="type.value",
              :label="type.label",
              color="primary"
            )
      template(#actions="")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS, ACK_MODAL_ACTIONS_TYPES } from '@/constants';

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
        output: '',
        ack_resources: false,
      },
    };
  },
  computed: {
    actionTypes() {
      return [
        {
          label: this.$t('alarm.acknowledge'),
          value: ACK_MODAL_ACTIONS_TYPES.ack,
        },
        {
          label: this.$t('alarm.acknowledgeAndAssociateTicket'),
          value: ACK_MODAL_ACTIONS_TYPES.ackAndAssociateTicket,
        },
        {
          label: this.$t('alarm.acknowledgeAndDeclareTicket'),
          value: ACK_MODAL_ACTIONS_TYPES.ackAndDeclareTicket,
        },
      ];
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
