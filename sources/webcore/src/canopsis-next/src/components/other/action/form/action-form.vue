<template lang="pug">
  v-form
    v-text-field(
    v-model="form.generalParameters._id",
    label="Id",
    )
    v-select(
    v-model="form.generalParameters.type",
    label="Type",
    :items="actionTypes"
    )
    v-tabs(centered, slider-color="primary")
      v-tab
        .validation-header(
        :class="{ 'error--text': hasGeneralFormAnyError }"
        ) {{ $t('modals.createAction.tabs.general') }}
      v-tab-item
        action-general-tab(v-model="form", ref="generalForm")
      v-tab
        .validation-header(
        :class="{ 'error--text': hasHookFormAnyError }"
        ) {{ $t('modals.createAction.tabs.hook') }}
      v-tab-item
        webhook-form-hook-tab(
        ref="hookForm",
        v-model="form.generalParameters.hook",
        :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS",
        )
</template>

<script>
import { ACTION_TYPES } from '@/constants';

import formMixin from '@/mixins/form';
import formValidationHeaderMixin from '@/mixins/form/validation-header';

import WebhookFormHookTab from '@/components/other/webhook/form/tabs/webhook-form-hook-tab.vue';
import ActionGeneralTab from './tabs/action-general-tab.vue';

export default {
  components: {
    ActionGeneralTab,
    WebhookFormHookTab,
  },
  mixins: [formMixin, formValidationHeaderMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      hasHookFormAnyError: false,
      hasGeneralFormAnyError: false,
    };
  },
  computed: {
    actionTypes() {
      return Object.values(ACTION_TYPES);
    },
  },
  mounted() {
    this.hasHookFormAnyError = false;
    this.hasGeneralFormAnyError = false;

    this.$watch(() => this.$refs.hookForm.hasAnyError, (value) => {
      this.hasHookFormAnyError = value;
    });

    this.$watch(() => this.$refs.generalForm.hasAnyError, (value) => {
      this.hasGeneralFormAnyError = value;
    });
  },
};
</script>
