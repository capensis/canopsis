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
    c-enabled-field(v-field="form.enabled")
    patterns-list(v-field="form.patterns")
</template>

<script>
import { EVENT_FILTER_RULE_TYPES } from '@/constants';

import formMixin from '@/mixins/form';

import PatternsList from '@/components/common/patterns-list/patterns-list.vue';

export default {
  inject: ['$validator'],
  components: { PatternsList },
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
};
</script>

