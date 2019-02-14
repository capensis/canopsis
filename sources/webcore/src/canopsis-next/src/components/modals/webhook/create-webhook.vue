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
              multiple,
              chips
              )
            v-flex(xs12)
              v-tabs(fixed-tabs)
                v-tab Events patterns
                v-tab Alarms patterns
                v-tab Entities patterns
                v-tab-item
                  v-layout(row, wrap)
                    v-flex(xs11)
                      v-textarea(
                      value="Test 123\nasd\nasdasd\nasdasd\nasd\nTEST",
                      rows="7",
                      no-resize,
                      readonly,
                      disabled
                      )
                    v-flex(xs1)
                      div
                        v-btn(icon, @click="showEditEventPatternModal")
                          v-icon edit
                      div
                        v-btn(color="error", icon, @click="showRemovePatternModal('event_pattern', 0)")
                          v-icon delete
                  v-btn(color="primary", @click="showEditEventPatternModal") Add pattern
                v-tab-item
                  v-btn(color="primary", @click="showEditEventPatternModal") Add pattern
                v-tab-item
                  v-btn(color="primary", @click="showEditEventPatternModal") Add pattern
        v-tab-item
          h2 Request
          v-layout(justify-space-between, align-center)
            v-flex
              v-text-field(
              v-model="form.request.method",
              label="Method",
              v-validate="'required'",
              name="request.method",
              :error-messages="errors.collect('request.method')"
              )
            v-flex
              v-text-field(
              v-model="form.request.url",
              label="URL",
              v-validate="'required'",
              name="request.url",
              :error-messages="errors.collect('request.url')"
              )
          v-layout
            v-flex
              v-text-field(
              label="Header key",
              v-validate="'required'",
              name="request.header.key",
              :error-messages="errors.collect('request.header.key')"
              )
            v-flex
              v-text-field(
              label="Header value",
              v-validate="'required'",
              name="request.header.value",
              :error-messages="errors.collect('request.header.value')"
              )
          v-layout
            v-flex
              v-textarea(
              v-model="form.request.payload",
              label="Payload",
              v-validate="'required'",
              name="request.payload",
              :error-messages="errors.collect('request.payload')"
              )
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS, WEBHOOK_TRIGGERS } from '@/constants';

import authMixin from '@/mixins/auth';
import popupMixin from '@/mixins/popup';
import modalInnerMixin from '@/mixins/modal/inner';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWebhook,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [
    authMixin,
    popupMixin,
    modalInnerMixin,
  ],
  data() {
    return {
      availableTriggers: Object.values(WEBHOOK_TRIGGERS),
      methods: [
        'POST', 'GET', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'CONNECT', 'OPTIONS', 'TRACE',
      ],
      form: {
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
          headers: {},
          payload: '{}',
        },
        declare_ticket: {},
      },
    };
  },
  methods: {
    showEditEventPatternModal() {
      this.showModal({
        name: MODALS.createEventFilterRulePattern,
        config: {},
      });
    },
    showRemovePatternModal(key, index) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.$delete(this.form.hook[key], index),
        },
      });
    },
    submit() {},
  },
};
</script>
