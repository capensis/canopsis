<template lang="pug">
  div
    v-text-field(
      v-field="form.generalParameters._id",
      v-validate="'required'",
      :error-messages="errors.collect('id')",
      :disabled="disabledId",
      label="Id",
      name="id"
    )
    v-select(
      v-field="form.generalParameters.type",
      v-validate="'required'",
      :items="actionTypes",
      :error-messages="errors.collect('actionType')",
      label="Type",
      name="actionType",
      @change="errors.clear()"
    )
    v-tabs(fixed-tabs, slider-color="primary")
      v-tab
        .validation-header(
          :class="{ 'error--text': hasGeneralFormAnyError }"
        ) {{ $t('modals.createAction.tabs.general') }}
      v-tab-item
        action-general-tab(
          ref="generalForm",
          v-field="form"
        )
      v-tab
        .validation-header(
          :class="{ 'error--text': hasHookFormAnyError }"
        ) {{ $t('modals.createAction.tabs.hook') }}
      v-tab-item
        webhook-form-hook-tab(
          ref="hookForm",
          v-field="form.generalParameters.hook"
        )
</template>

<script>
import { ACTION_TYPES } from '@/constants';

import formValidationHeaderMixin from '@/mixins/form/validation-header';

import WebhookFormHookTab from '@/components/other/webhook/form/tabs/webhook-form-hook-tab.vue';
import ActionGeneralTab from './tabs/action-general-tab.vue';

export default {
  inject: ['$validator'],
  components: {
    ActionGeneralTab,
    WebhookFormHookTab,
  },
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    disabledId: {
      type: Boolean,
      default: false,
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
      return Object.values(ACTION_TYPES).map((type) => {
        let text = type;

        if (type === ACTION_TYPES.changeState) {
          text = `${type} (${this.$t('alarmList.actions.titles.changeState')})`;
        }

        return { text, value: type };
      });
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
