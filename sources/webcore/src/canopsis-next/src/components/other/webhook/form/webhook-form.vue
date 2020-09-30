<template lang="pug">
  v-tabs(
    :color="vTabsColor",
    :dark="dark",
    fixed-tabs,
    slider-color="primary"
  )
    template(v-for="tab in tabs")
      v-tab(:key="`tab-${tab.key}`")
        div.validation-header(:class="{ 'error--text': validationErrorsFlagsForTabs[tab.key] }") {{ tab.title }}
      v-tab-item(:key="`tab-item-${tab.key}`")
        div(:class="vTabItemInnerWrapperClass")
          div(:class="vTabItemInnerClass")
            component.pt-2(
              :ref="tab.key",
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
      validationErrorsFlagsForTabs: {},
      watchers: [],
    };
  },
  computed: {
    tabs() {
      const tabs = [
        {
          key: 'hook',
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
          key: 'request',
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
          key: 'declare-ticket',
          title: this.$t('webhook.tabs.declareTicket.title'),
          component: 'webhook-form-declare-ticket-tab',
          bind: {
            declareTicket: this.form.declare_ticket,
            disabled: this.disabled,
            emptyResponse: this.form.emptyResponse,
          },
          on: {
            input: event => this.updateField('declare_ticket', event),
            'update:emptyResponse': event => this.updateField('emptyResponse', event),
          },
        },
      ];

      if (this.form.combine_meta_alarm_request) {
        tabs.push({
          key: 'combine-meta-alarm-request',
          title: this.$t('webhook.tabs.combineMetaAlarmRequest.title'),
          component: 'webhook-form-request-tab',
          bind: {
            request: this.form.combine_meta_alarm_request,
            disabled: this.disabled,
            namePrefix: 'combineMetaAlarmRequest',
          },
          on: {
            input: event => this.updateField('combine_meta_alarm_request', event),
          },
        });
      }

      return tabs;
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
  watch: {
    'form.combine_meta_alarm_request': {
      handler() {
        this.$nextTick(() => {
          this.setupWatchers();
        });
      },
    },
  },
  mounted() {
    this.setupWatchers();
  },
  methods: {
    setupWatchers() {
      const formRefs = this.tabs.map(({ key }) => key);

      if (this.watchers.length) {
        this.watchers.forEach(unwatch => unwatch());
      }

      this.watchers = formRefs.reduce((acc, key) => {
        if (!this.$refs[key]) {
          return acc;
        }

        const [ref] = this.$refs[key];

        if (ref) {
          this.$set(this.validationErrorsFlagsForTabs, key, ref.hasAnyError);

          acc.push(this.$watch(() => ref.hasAnyError, (value) => {
            this.$set(this.validationErrorsFlagsForTabs, key, value);
          }));
        }

        return acc;
      }, []);
    },
  },
};
</script>
