<template lang="pug">
  v-layout(column)
    v-flex(v-if="!webhooks.length", xs12)
      v-alert(:value="true", type="info") {{ $t('declareTicket.emptyWebhooks') }}
    draggable(v-field="webhooks", :options="draggableOptions")
      v-card.my-2(v-for="(webhook, index) in webhooks", :key="webhook.key")
        v-card-text
          declare-ticket-rule-webhook-field(v-field="webhooks[index]", :name="webhook.key")
    v-layout(row, align-center)
      v-btn.ml-0(
        :color="hasWebhooksErrors ? 'error' : 'primary'",
        :disabled="disabled",
        outline,
        @click="addWebhook"
      ) {{ $t('declareTicket.addWebhook') }}
      span.error--text(v-show="hasWebhooksErrors") {{ $t('declareTicket.errors.webhookRequired') }}
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { declareTicketRuleWebhookToForm } from '@/helpers/forms/declare-ticket-rule';

import { formArrayMixin } from '@/mixins/form';

import DeclareTicketRuleWebhookField from './declare-ticket-rule-webhook-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DeclareTicketRuleWebhookField,
    Draggable,
  },
  mixins: [formArrayMixin],
  model: {
    prop: 'webhooks',
    event: 'input',
  },
  props: {
    webhooks: {
      type: Array,
      default: () => ([]),
    },
    name: {
      type: String,
      default: 'webhooks',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    hasWebhooksErrors() {
      return this.errors.has(this.name);
    },

    draggableOptions() {
      return {
        disabled: this.disabled,
        animation: VUETIFY_ANIMATION_DELAY,
        handle: '.step-drag-handler',
        ghostClass: 'white',
        group: {
          name: 'remediation-instruction-steps',
        },
      };
    },
  },
  watch: {
    webhooks() {
      this.$validator.validate(this.name);
    },
  },
  created() {
    this.$validator.attach({
      name: this.name,
      rules: 'min_value:1',
      getter: () => this.webhooks.length,
      vm: this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.name);
  },
  methods: {
    addWebhook() {
      this.addItemIntoArray(declareTicketRuleWebhookToForm());
    },

    removeStep(index) {
      this.removeItemFromArray(index);
    },
  },
};
</script>
