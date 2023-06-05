<template lang="pug">
  v-tabs(color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $tc('common.pattern', 2) }}
    v-tab-item(lazy)
      v-layout.py-3
        v-flex(xs12, md8, offset-md2)
          v-card
            v-card-text
              meta-alarm-rule-patterns-form(
                :form="patterns",
                :with-total-entity="withTotalEntityPattern",
                readonly
              )
</template>

<script>
import {
  isMetaAlarmRuleTypeHasTotalEntityPatterns,
  metaAlarmFilterPatternsToForm,
} from '@/helpers/forms/meta-alarm-rule';

import MetaAlarmRulePatternsForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-patterns-form.vue';

export default {
  components: { MetaAlarmRulePatternsForm },
  props: {
    metaAlarmRule: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    patterns() {
      return metaAlarmFilterPatternsToForm(this.metaAlarmRule);
    },

    withTotalEntityPattern() {
      return isMetaAlarmRuleTypeHasTotalEntityPatterns(this.metaAlarmRule.type);
    },
  },
};
</script>
