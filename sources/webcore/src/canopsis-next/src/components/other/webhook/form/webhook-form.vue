<template lang="pug">
  v-tabs(:color="vTabsColor", :dark="dark", fixed-tabs)
    template(v-for="(tab, index) in tabs")
      v-tab(:key="`tab-${index}`")
        .validation-header(:class="{ 'error--text': validationErrorsFlagsForTabs[index] }") {{ tab.title }}
      v-tab-item(:key="`tab-item-${index}`")
        div(:class="vTabItemInnerWrapperClass")
          div(:class="vTabItemInnerClass")
            component(
            ref="forms",
            :is="tab.component",
            :class="webhookTabClass",
            v-bind="tab.bind",
            v-on="tab.on"
            )
</template>

<script>
import { intersection } from 'lodash';

import { WEBHOOK_TRIGGERS } from '@/constants';

import formMixin from '@/mixins/form';

import WebhookFormHookTab from './tabs/webhook-form-hook-tab.vue';
import WebhookFormRequestTab from './tabs/webhook-form-request-tab.vue';
import WebhookFormDeclareTicketTab from './tabs/webhook-form-declare-ticket-tab.vue';

export default {
  components: {
    WebhookFormHookTab,
    WebhookFormRequestTab,
    WebhookFormDeclareTicketTab,
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
    color: {
      type: String,
      default: null,
    },
    dark: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      validationErrorsFlagsForTabs: [],
    };
  },
  computed: {
    tabs() {
      return [
        {
          title: this.$t('webhook.tabs.hook.title'),
          component: 'webhook-form-hook-tab',
          bind: {
            hook: this.form.hook,
            hasBlockedTriggers: this.hasBlockedTriggers,
            disabled: this.disabled,
          },
          on: {
            input: event => this.updateField('hook', event),
          },
        },
        {
          title: this.$t('webhook.tabs.request.title'),
          component: 'webhook-form-request-tab',
          bind: {
            request: this.form.request,
            disabled: this.disabled,
          },
          on: {
            input: event => this.updateField('request', event),
          },
        },
        {
          title: this.$t('webhook.tabs.declareTicket.title'),
          component: 'webhook-form-declare-ticket-tab',
          bind: {
            declareTicket: this.form.declare_ticket,
            disabled: this.disabled,
          },
          on: {
            input: event => this.updateField('declare_ticket', event),
          },
        },
      ];
    },

    hasBlockedTriggers() {
      return intersection(this.form.hook.triggers, [
        WEBHOOK_TRIGGERS.resolve,
        WEBHOOK_TRIGGERS.unsnooze,
      ]).length > 0;
    },

    vTabsColor() {
      return this.dark ? 'secondary lighten-1' : null;
    },

    vTabItemInnerWrapperClass() {
      return {
        'secondary lighten-2': this.dark,
      };
    },

    vTabItemInnerClass() {
      return {
        'pa-3': this.dark,
      };
    },

    webhookTabClass() {
      return {
        'white pa-3': this.dark,
      };
    },
  },
  mounted() {
    this.tabs.forEach((item, index) => {
      this.$set(this.validationErrorsFlagsForTabs, index, false);

      this.$watch(() => this.$refs.forms[index].hasAnyError, (value) => {
        this.$set(this.validationErrorsFlagsForTabs, index, value);
      });
    });
  },
};
</script>
