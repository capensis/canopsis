<template lang="pug">
  div
    v-layout(row, wrap)
      v-flex(xs12)
        v-select(
          v-field="hook.triggers",
          v-validate="'required'",
          :items="availableTriggers",
          :disabled="disabled",
          :label="$t('webhook.tabs.hook.fields.triggers')",
          :error-messages="errors.collect('triggers')",
          name="triggers",
          multiple,
          chips
        )
      v-flex(xs12)
        v-tabs(v-model="activeHookTab", fixed-tabs)
          v-tab(:disabled="hasBlockedTriggers") {{ $t('webhook.tabs.hook.fields.eventPatterns') }}
          v-tab {{ $t('webhook.tabs.hook.fields.alarmPatterns') }}
          v-tab {{ $t('webhook.tabs.hook.fields.entityPatterns') }}
          v-tab-item(:disabled="hasBlockedTriggers")
            patterns-list(
              v-field="hook.event_patterns",
              :disabled="disabled",
              :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS"
            )
          v-tab-item
            patterns-list(
              v-field="hook.alarm_patterns",
              :disabled="disabled",
              :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS"
            )
          v-tab-item
            patterns-list(
              v-field="hook.entity_patterns",
              :disabled="disabled",
              :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS"
            )
</template>

<script>
import { WEBHOOK_TRIGGERS } from '@/constants';

import formValidationHeaderMixin from '@/mixins/form/validation-header';

import PatternsList from '@/components/other/shared/patterns-list/patterns-list.vue';

export default {
  inject: ['$validator'],
  components: { PatternsList },
  mixins: [formValidationHeaderMixin],
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
    disabled: {
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
