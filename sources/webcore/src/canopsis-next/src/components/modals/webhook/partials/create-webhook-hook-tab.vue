<template lang="pug">
  div
    h2 {{ $t('webhook.tabs.hook.title') }}
    v-layout(row, wrap)
      v-flex(xs12)
        v-select(
        :value="hook.triggers",
        :items="availableTriggers",
        :label="$t('webhook.tabs.hook.fields.triggers')",
        :error-messages="errors.collect('triggers')",
        v-validate="'required'",
        name="triggers",
        multiple,
        chips,
        @input="updateField('triggers', $event)"
        )
      v-flex(xs12)
        v-tabs(v-model="activeHookTab", fixed-tabs)
          v-tab(:disabled="hasBlockedTriggers") {{ $t('webhook.tabs.hook.fields.eventPatterns') }}
          v-tab {{ $t('webhook.tabs.hook.fields.alarmPatterns') }}
          v-tab {{ $t('webhook.tabs.hook.fields.entityPatterns') }}
          v-tab-item(:disabled="hasBlockedTriggers")
            webhook-hook-tab-patterns-list(
            :patterns="hook.event_pattern",
            @input="updateField('event_pattern', $event)"
            )
          v-tab-item
            webhook-hook-tab-patterns-list(
            :patterns="hook.alarm_pattern",
            @input="updateField('alarm_pattern', $event)"
            )
          v-tab-item
            webhook-hook-tab-patterns-list(
            :patterns="hook.entity_pattern",
            @input="updateField('entity_pattern', $event)"
            )
</template>

<script>
import { WEBHOOK_TRIGGERS } from '@/constants';

import formMixin from '@/mixins/form';

import WebhookHookTabPatternsList from './create-webhook-hook-tab-patterns-list.vue';

export default {
  inject: ['$validator'],
  components: { WebhookHookTabPatternsList },
  mixins: [formMixin],
  model: {
    prop: 'hook',
    event: 'input',
  },
  props: {
    hook: {
      type: Object,
      required: true,
    },
    hasBlockedTriggers: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      activeHookTab: 0,
      availableTriggers: Object.values(WEBHOOK_TRIGGERS),
    };
  },
  watch: {
    hasBlockedTriggers: {
      immediate: true,
      handler(value, oldValue) {
        if (value && value !== oldValue && this.activeHookTab === 0) {
          this.activeHookTab = 1;
        }
      },
    },
  },
};
</script>
