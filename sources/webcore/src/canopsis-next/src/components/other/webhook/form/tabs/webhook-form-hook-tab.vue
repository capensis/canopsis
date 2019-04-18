<template lang="pug">
  div
    v-layout(row, wrap)
      v-flex(xs12)
        v-select(
        :value="hook.triggers",
        :items="availableTriggers",
        :disabled="disabled",
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
            patterns-list(
            :patterns="hook.event_patterns",
            :disabled="disabled",
            :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS",
            @input="updateField('event_patterns', $event)"
            )
          v-tab-item
            patterns-list(
            :patterns="hook.alarm_patterns",
            :disabled="disabled",
            :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS",
            @input="updateField('alarm_patterns', $event)"
            )
          v-tab-item
            patterns-list(
            :patterns="hook.entity_patterns",
            :disabled="disabled",
            :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS",
            @input="updateField('entity_patterns', $event)"
            )
</template>

<script>
import { WEBHOOK_TRIGGERS } from '@/constants';

import formMixin from '@/mixins/form';

import PatternsList from '@/components/other/shared/patterns-list/patterns-list.vue';

export default {
  inject: ['$validator'],
  components: { PatternsList },
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
