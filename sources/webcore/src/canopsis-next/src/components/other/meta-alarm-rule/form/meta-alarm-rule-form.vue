<template lang="pug">
  div
    v-layout(align-center)
      v-text-field(
        v-field="form._id",
        :label="$t('metaAlarmRule.id')",
        :disabled="isDisabledIdField",
        :readonly="isDisabledIdField"
      )
        v-tooltip(v-show="!isDisabledIdField", slot="append", left)
          v-icon(slot="activator") help
          span {{ $t('metaAlarmRule.idHelp') }}
    v-select(v-field="form.type", :items="ruleTypes", :label="$t('common.type')")
    v-textarea(v-field="form.description", :label="$t('common.description')")
    v-btn(@click="editPattern") {{ $t('modals.metaAlarmRule.editPattern') }}
</template>

<script>
import { MODALS, META_ALARMS_RULE_TYPES } from '@/constants';

import formMixin from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    ruleTypes() {
      return Object.values(META_ALARMS_RULE_TYPES);
    },
  },
  methods: {
    editPattern() {
      this.$modals.show({
        name: MODALS.createEventFilterRulePattern,
        config: {
          pattern: this.form.alarm_patterns,
          action: alarmPatterns => this.updateField('alarm_patterns', alarmPatterns),
        },
      });
    },
  },
};
</script>

