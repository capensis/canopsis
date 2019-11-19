<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createAckEvent.title') }}
      template(slot="text")
        v-container
          v-layout(row, align-center)
            v-flex.text-xs-center
              alarm-general-table(:items="items")
          v-layout(row)
            v-divider.my-3
          v-layout(row)
            v-text-field(
              v-model="form.ticket",
              :label="$t('modals.createAckEvent.fields.ticket')"
            )
          v-layout(row)
            v-textarea(
              v-model="form.output",
              v-validate="isNoteRequired ? 'required' : ''",
              :label="$t('modals.createAckEvent.fields.output')",
              :error-messages="errors.collect('output')",
              name="output"
            )
          v-layout(row)
            v-tooltip(top)
              v-checkbox(
                slot="activator",
                v-model="ack_resources",
                :label="$t('modals.createAckEvent.fields.ackResources')"
              )
              span {{ $t('modals.createAckEvent.tooltips.ackResources') }}
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="errors.any() || submitting",
          type="submit"
        ) {{ $t('common.actions.ack') }}
        v-btn.warning(
          :disabled="errors.any() || submitting",
          @click="submitWithTicket"
        ) {{ submitWithTicketBtnLabel }}
</template>

<script>
import { omit } from 'lodash';
import { MODALS, EVENT_ENTITY_TYPES } from '@/constants';

import AlarmGeneralTable from '@/components/other/alarm/alarm-general-list.vue';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create an ack event
 */
export default {
  name: MODALS.createAckEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AlarmGeneralTable, ModalWrapper },
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin],
  data() {
    return {
      submitting: false,
      ack_resources: false,
      form: {
        ticket: '',
        output: '',
      },
    };
  },
  computed: {
    isNoteRequired() {
      return this.config && this.config.isNoteRequired;
    },

    submitWithTicketBtnLabel() {
      return this.form.ticket ? this.$t('common.actions.acknowledgeAndAssociateTicket') : this.$t('common.actions.acknowledgeAndDeclareTicket');
    },
  },
  methods: {
    createAckEvent() {
      const ackEventData = this.prepareData(EVENT_ENTITY_TYPES.ack, this.items, this.form);

      return this.createEventAction({ data: ackEventData });
    },

    createDeclareTicketEvent() {
      const declareTicketEventData = this.prepareData(EVENT_ENTITY_TYPES.declareTicket, this.items, omit(this.form, ['ticket']));

      return this.createEventAction({ data: declareTicketEventData });
    },

    createAssocTicketEvent() {
      const assocTicketEventData = this.prepareData(EVENT_ENTITY_TYPES.assocTicket, this.items, this.form);

      return this.createEventAction({ data: assocTicketEventData });
    },

    async createAckEventAndCloseModal() {
      await this.createAckEvent();

      if (this.config && this.config.afterSubmit) {
        await this.config.afterSubmit();
      }

      this.$modals.hide();
    },

    async submitWithTicket() {
      const formIsValid = await this.$validator.validateAll();

      if (formIsValid) {
        if (this.form.ticket) {
          await this.createAssocTicketEvent();
        } else {
          await this.createDeclareTicketEvent();
        }

        this.createAckEventAndCloseModal();
      }
    },

    async submit() {
      try {
        this.submitting = true;

        const formIsValid = await this.$validator.validateAll();

        if (formIsValid) {
          if (this.form.ticket) {
            this.$modals.show({
              name: MODALS.confirmAckWithTicket,
              config: {
                continueAction: this.createAckEventAndCloseModal,
                continueWithTicketAction: this.submitWithTicket,
              },
            });
          } else {
            await this.createAckEventAndCloseModal();
          }
        }
      } catch (err) {
        this.$popups.error({ text: err.description || this.$t('error.default') });
      } finally {
        this.submitting = false;
      }
    },
  },
};
</script>

