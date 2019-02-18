<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Title
    v-card-text
      v-tabs(fixed-tabs)
        v-tab Hook
        v-tab Request
        v-tab Declare ticket
        v-tab-item
          h2 Hook
          v-layout(row, wrap)
            v-flex(xs12)
              v-select(
              v-model="form.hook.triggers",
              :items="availableTriggers",
              label="Triggers",
              :error-messages="errors.collect('triggers')",
              v-validate="'required'",
              name="triggers",
              multiple,
              chips
              )
            v-flex(xs12)
              v-tabs(fixed-tabs)
                v-tab Events patterns
                v-tab Alarms patterns
                v-tab Entities patterns
                v-tab-item
                  patterns-list(v-model="form.hook.event_pattern")
                v-tab-item
                  patterns-list(v-model="form.hook.alarm_pattern")
                v-tab-item
                  patterns-list(v-model="form.hook.entity_pattern")
        v-tab-item
          webhook-request(v-model="form.request")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS, WEBHOOK_TRIGGERS } from '@/constants';

import uid from '@/helpers/uid';
import { setInWith } from '@/helpers/immutable';

import authMixin from '@/mixins/auth';
import popupMixin from '@/mixins/popup';
import modalInnerMixin from '@/mixins/modal/inner';

import PatternsList from './partials/patterns-list.vue';
import WebhookRequest from './partials/webhook-request.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWebhook,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    PatternsList,
    WebhookRequest,
  },
  mixins: [
    authMixin,
    popupMixin,
    modalInnerMixin,
  ],
  data() {
    const { webhook } = this.modal.config;
    const defaultForm = {
      hook: {
        triggers: [],
        event_pattern: [],
        alarm_pattern: [],
        entity_pattern: [],
      },
      request: {
        method: '',
        url: '',
        auth: {
          username: '',
          password: '',
        },
        headers: [],
        payload: '{}',
      },
      declare_ticket: {},
    };

    return {
      availableTriggers: Object.values(WEBHOOK_TRIGGERS),
      methods: [
        'POST', 'GET', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'CONNECT', 'OPTIONS', 'TRACE',
      ],
      form: webhook ? this.prepareWebhookToForm(webhook) : defaultForm,
    };
  },
  methods: {
    prepareWebhookToForm(webhook) {
      return setInWith(webhook, 'request.headers', headers =>
        Object.keys(headers).map(key => ({ key, id: uid(), value: headers[key] })));
    },

    prepareWebhookToRequest(webhook) {
      return setInWith(webhook, 'request.headers', headers =>
        headers.reduce((acc, header) => {
          acc[header.key] = header.value;

          return acc;
        }, {}));
    },

    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.hideModal();
      }
    },
  },
};
</script>
