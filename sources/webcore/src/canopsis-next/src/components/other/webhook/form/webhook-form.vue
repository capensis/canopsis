<template lang="pug">
  div
    v-text-field(
      v-field="form._id",
      :label="$t('modals.createWebhook.fields.id')",
      :readonly="isDisabledIdField",
      :disabled="isDisabledIdField"
    )
      v-tooltip(v-show="!isDisabledIdField", slot="append", left)
        v-icon(slot="activator") help_outline
        span {{ $t('modals.createWebhook.tooltips.id') }}
    retry-field(v-field="form.retry")
    c-enabled-field(v-field="form.enabled")
    c-disable-during-periods-field(v-field="form.disable_during_periods")
    webhook-tabs-form(v-field="form")
</template>

<script>
import formMixin from '@/mixins/form';

import RetryField from '@/components/forms/fields/retry.vue';

import WebhookTabsForm from './tabs/webhook-tabs-form.vue';

export default {
  components: {
    RetryField,
    WebhookTabsForm,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
