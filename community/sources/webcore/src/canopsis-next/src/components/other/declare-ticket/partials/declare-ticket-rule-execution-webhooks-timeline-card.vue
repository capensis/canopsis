<template>
  <v-layout column>
    <v-layout
      class="mb-2"
      align-center
    >
      <c-expand-btn
        v-if="expandable"
        :expanded="expanded"
        :loading="pending"
        color="grey"
        @expand="expandResponse"
      />
      <span class="ml-2">{{ $t(`alarm.timeLine.types.${step._t}`) }}</span>
      <span v-if="step.message">: {{ step.message }}</span>
    </v-layout>
    <v-expand-transition>
      <v-card
        v-show="expanded"
        light
      >
        <v-card-text>
          <c-request-text-information
            v-if="response"
            :value="response"
          />
          <span v-else>{{ $t('declareTicket.noResponse') }}</span>
        </v-card-text>
      </v-card>
    </v-expand-transition>
  </v-layout>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { EVENT_ENTITY_TYPES } from '@/constants';

const { mapActions: mapDeclareTicketRuleActions } = createNamespacedHelpers('declareTicketRule');

export default {
  props: {
    step: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      expanded: false,
      response: null,
    };
  },
  computed: {
    expandable() {
      return [EVENT_ENTITY_TYPES.webhookComplete, EVENT_ENTITY_TYPES.webhookFail].includes(this.step._t);
    },
  },
  methods: {
    ...mapDeclareTicketRuleActions({
      fetchTestDeclareTicketExecutionWebhooksResponse: 'fetchTestDeclareTicketExecutionWebhooksResponse',
    }),

    async expandResponse(value) {
      if (value) {
        if (!this.response) {
          this.pending = true;

          try {
            this.response = await this.fetchTestDeclareTicketExecutionWebhooksResponse({ id: this.step.webhookId });
          } catch (err) {
            console.error(err);
          } finally {
            this.pending = false;
          }
        }

        this.expanded = true;
      } else {
        this.expanded = false;
      }
    },
  },
};
</script>
