<template lang="pug">
  div
    v-layout(align-center)
      v-text-field(
        v-field="form._id",
        :label="$t('eventFilter.id')",
        :disabled="isDisabledIdField",
        :readonly="isDisabledIdField"
      )
        v-tooltip(v-show="!isDisabledIdField", slot="append", left)
          v-icon(slot="activator") help
          span {{ $t('eventFilter.idHelp') }}
    v-select(
      v-field="form.type",
      :items="ruleTypes",
      :label="$t('common.type')"
    )
    v-textarea(v-field="form.description", :label="$t('common.description')")
    v-text-field(
      v-field.number="form.priority",
      :label="$t('modals.eventFilterRule.priority')",
      type="number"
    )
    v-switch(v-field="form.enabled", :label="$t('common.enabled')")
    v-btn(@click="editPattern") {{ $t('modals.eventFilterRule.editPattern') }}
</template>

<script>
import { MODALS, EVENT_FILTER_RULE_TYPES } from '@/constants';

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
      return Object.values(EVENT_FILTER_RULE_TYPES);
    },
  },
  methods: {
    editPattern() {
      this.$modals.show({
        name: MODALS.createEventFilterRulePattern,
        config: {
          pattern: this.form.pattern,
          action: pattern => this.updateField('pattern', pattern),
        },
      });
    },
  },
};
</script>

