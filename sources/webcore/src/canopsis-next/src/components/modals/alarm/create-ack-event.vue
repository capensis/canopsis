<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createAckEvent.title') }}
    v-card-text
      v-container
        v-layout(row, align-center)
          v-flex.text-xs-center
            alarm-general-table(:items="items")
        v-layout(row)
          v-divider.my-3
        v-layout(row)
          v-text-field(
            :label="$t('modals.createAckEvent.fields.ticket')",
            v-model="form.ticket"
          )
        v-layout(row)
          v-textarea(
            :label="$t('modals.createAckEvent.fields.output')",
            :error-messages="errors.collect('output')",
            v-model="form.output",
            v-validate="isNoteRequired ? 'required' : ''",
            data-vv-name="output"
          )
        v-layout(row)
          v-tooltip(top)
            v-checkbox(
              slot="activator",
              v-model="ack_resources",
              :label="$t('modals.createAckEvent.fields.ackResources')"
            )
              span(slot-name="label") {{  }}
            span {{ $t('modals.createAckEvent.tooltips.ackResources') }}
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.actions.ack') }}
      v-btn.warning(@click.prevent="submitWithTicket") {{ submitWithTicketBtnLabel }}
</template>

<script>
import { omit } from 'lodash';
import { MODALS, EVENT_ENTITY_TYPES } from '@/constants';

import AlarmGeneralTable from '@/components/other/alarm/alarm-general-list.vue';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';

/**
 * Modal to create an ack event
 */
export default {
  name: MODALS.createAckEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    AlarmGeneralTable,
  },
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin],
  data() {
    return {
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
          this.createAckEventAndCloseModal();
        }
      }
    },
  },
};
</script>

<style scoped>
  .tooltip {
    flex: 1 1 auto;
  }
</style>
