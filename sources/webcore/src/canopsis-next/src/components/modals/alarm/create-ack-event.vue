<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createAckEvent.title') }}
    v-card-text
      v-container
        v-layout(row align-center)
          v-flex.text-xs-center
            alarm-general-table(:items="items")
        v-layout(row)
          v-divider.my-3
        v-layout(row)
          v-text-field(
          :label="$t('modals.createAckEvent.fields.ticket')",
          :error-messages="errors.collect('ticket')",
          v-model="form.ticket",
          v-validate="rules",
          data-vv-name="ticket"
          )
        v-layout(row)
          v-textarea(
          :label="$t('modals.createAckEvent.fields.output')",
          :error-messages="errors.collect('output')",
          v-model="form.output",
          v-validate="rules",
          data-vv-name="output",
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
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit", :disabled="errors.any()") {{ $t('common.actions.ack') }}
      v-btn.warning(
      @click.prevent="submitWithDeclare",
      ) {{ $t('common.actions.acknowledgeAndReport') }}
</template>

<script>
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
      showValidationErrors: false,
      ack_resources: false,
      form: {
        ticket: '',
        output: '',
      },
    };
  },
  computed: {
    rules() {
      return this.showValidationErrors ? 'required' : '';
    },
  },
  methods: {
    async create(withDeclare) {
      const ackEventData = this.prepareData(EVENT_ENTITY_TYPES.ack, this.items, this.form);

      await this.createEventAction({
        data: ackEventData,
      });

      if (withDeclare) {
        const declareTicketEventData =
          this.prepareData(EVENT_ENTITY_TYPES.declareTicket, this.items, this.form);

        await this.createEventAction({
          data: declareTicketEventData,
        });
      }

      if (this.config && this.config.afterSubmit) {
        await this.config.afterSubmit();
      }

      this.hideModal();
    },

    async submitWithDeclare() {
      const formIsValid = await this.$validator.validateAll();

      if (formIsValid) {
        await this.create(true);
      }
    },

    async submit() {
      this.showValidationErrors = false;
      this.errors.clear();

      await this.create();
    },
  },
};
</script>

<style scoped>
  .tooltip {
    flex: 1 1 auto;
  }
</style>
