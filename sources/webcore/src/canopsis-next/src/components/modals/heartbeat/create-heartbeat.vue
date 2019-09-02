<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ title }}
    v-card-text
      v-text-field(
        v-model="form._id",
        :label="$t('modals.createWebhook.fields.id')",
        :readonly="isDisabledIdField",
        :disabled="isDisabledIdField"
      )
        v-tooltip(slot="append", left)
          v-icon(slot="activator") help_outline
          span {{ $t('modals.createWebhook.tooltips.id') }}
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="hideModal") {{ $t('common.cancel') }}
      v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createHeartbeat,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin],
  data() {
    const { webhook } = this.modal.config;
    const defaultForm = {
      hook: {
        triggers: [],
        event_patterns: [],
        alarm_patterns: [],
        entity_patterns: [],
      },
      request: {
        method: '',
        url: '',
        headers: [],
        payload: '{}',
      },
      declare_ticket: [],
      disable_if_active_pbehavior: false,
      emptyResponse: false,
    };

    return {
      form: webhook ? this.$options.filters.webhookToForm(cloneDeep(webhook)) : defaultForm,
    };
  },
  computed: {
    title() {
      if (this.config.heartbeat) {
        return this.$t('modals.createHeartbeat.edit.title');
      }

      return this.$t('modals.createHeartbeat.create.title');
    },
    isDisabledIdField() {
      return !!this.config.webhook;
    },
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        this.hideModal();
      }
    },
  },
};
</script>
