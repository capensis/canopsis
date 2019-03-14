<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ title }}
    v-card-text
      webhook-form(v-model="form")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import { setInSeveral } from '@/helpers/immutable';

import modalInnerMixin from '@/mixins/modal/inner';
import webhookFormFiltersMixin from '@/mixins/webhook/form-filters';

import WebhookForm from '@/components/other/webhook/form/webhook-form.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWebhook,
  $_veeValidate: {
    validator: 'new',
  },
  components: { WebhookForm },
  mixins: [modalInnerMixin, webhookFormFiltersMixin],
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
        auth: {
          username: '',
          password: '',
        },
        headers: [],
        payload: '{}',
      },
      declare_ticket: [],
    };

    return {
      form: webhook ? this.$options.filters.webhookToForm(cloneDeep(webhook)) : defaultForm,
    };
  },
  computed: {
    title() {
      if (this.config.webhook) {
        return this.$t('modals.createWebhook.edit.title');
      }

      return this.$t('modals.createWebhook.create.title');
    },
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        if (this.config.action) {
          const data = {
            ...this.form,
            hook: {
              ...this.form.hook,
              event_patterns: this.form.hook.event_patterns.length ? this.form.hook.event_patterns : null,
              alarm_patterns: this.form.hook.alarm_patterns.length ? this.form.hook.alarm_patterns : null,
              entity_patterns: this.form.hook.entity_patterns.length ? this.form.hook.entity_patterns : null,
            },
          };

          const preparedForm = this.hasBlockedTriggers ? setInSeveral(data, {
            'hook.event_patterns': null,
            declare_ticket: {},
          }) : data;

          await this.config.action(this.$options.filters.formToWebhook(preparedForm));
        }

        this.hideModal();
      }
    },
  },
};
</script>
