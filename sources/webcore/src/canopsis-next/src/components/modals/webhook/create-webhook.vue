<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        v-text-field(
          v-model="form._id",
          :label="$t('modals.createWebhook.fields.id')",
          :readonly="isDisabledIdField",
          :disabled="isDisabledIdField"
        )
          v-tooltip(v-show="!isDisabledIdField", slot="append", left)
            v-icon(slot="activator") help_outline
            span {{ $t('modals.createWebhook.tooltips.id') }}
        retry-field(v-model="form.retry")
        enabled-field(v-model="form.enabled")
        disable-during-periods-field(v-model="form.disable_during_periods")
        webhook-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { omit } from 'lodash';

import { MODALS } from '@/constants';

import { setSeveralFields } from '@/helpers/immutable';
import { formToWebhook, webhookToForm } from '@/helpers/forms/webhook';

import authMixin from '@/mixins/auth';
import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import DisableDuringPeriodsField from '@/components/forms/fields/disable-during-periods.vue';
import WebhookForm from '@/components/other/webhook/form/webhook-form.vue';
import RetryField from '@/components/forms/fields/retry.vue';
import EnabledField from '@/components/forms/fields/enabled-field.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWebhook,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    EnabledField,
    DisableDuringPeriodsField,
    WebhookForm,
    ModalWrapper,
    RetryField,
  },
  mixins: [authMixin, modalInnerMixin, submittableMixin()],
  data() {
    const preparedWebhook = this.modal.config.isDuplicating
      ? omit(this.modal.config.webhook, ['_id'])
      : this.modal.config.webhook;

    return {
      form: webhookToForm(preparedWebhook),
    };
  },
  computed: {
    title() {
      let type = 'create';

      if (this.config.webhook) {
        type = this.config.isDuplicating ? 'duplicate' : 'edit';
      }

      return this.$t(`modals.createWebhook.${type}.title`);
    },

    isDisabledIdField() {
      return this.config.webhook && !this.config.isDuplicating;
    },
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        if (this.config.action) {
          const preparedForm = this.hasBlockedTriggers ? setSeveralFields(this.form, {
            'hook.event_patterns': null,
            declare_ticket: {},
          }) : this.form;

          const webhook = formToWebhook(preparedForm);
          webhook.author = this.currentUser._id;

          await this.config.action(webhook);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
