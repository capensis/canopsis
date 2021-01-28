<template lang="pug">
  v-card
    v-card-text
      v-layout(row)
        c-action-type-field(v-field="action.type", :disabled="disabled")
      v-layout(row)
        c-enabled-field(v-field="action.emit_trigger", :label="$t('scenario.fields.emitTrigger')")
      v-tabs(centered, slider-color="primary")
        v-tab {{ $t('scenario.tabs.general') }}
        v-tab-item
          v-divider
          scenario-action-general-field(v-field="action")
        v-tab {{ $t('scenario.tabs.pattern') }}
        v-tab-item
          v-divider
          c-patterns-field(v-model="patterns", alarm, entity)
</template>

<script>
import formMixin from '@/mixins/form/object';

import ScenarioActionGeneralField from './scenario-action-general-field.vue';

export default {
  components: { ScenarioActionGeneralField },
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'action',
    event: 'input',
  },
  props: {
    action: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    patterns: {
      get() {
        return {
          alarm_patterns: this.action.alarm_patterns,
          entity_patterns: this.action.entity_patterns,
        };
      },
      set(patterns) {
        this.updateModel({
          ...this.action,
          ...patterns,
        });
      },
    },
  },
};
</script>
