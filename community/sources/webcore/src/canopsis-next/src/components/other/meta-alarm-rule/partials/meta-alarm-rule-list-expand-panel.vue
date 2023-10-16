<template>
  <v-tabs
    color="secondary lighten-1"
    slider-color="primary"
    dark="dark"
    centered="centered"
  >
    <v-tab>{{ $tc('common.pattern', 2) }}</v-tab>
    <v-tab-item>
      <v-layout class="py-3">
        <v-flex
          xs12="xs12"
          md8="md8"
          offset-md2="offset-md2"
        >
          <v-card>
            <v-card-text>
              <meta-alarm-rule-patterns-form
                :form="patterns"
                :with-total-entity="withTotalEntityPattern"
                readonly="readonly"
              />
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>
  </v-tabs>
</template>

<script>
import {
  isMetaAlarmRuleTypeHasTotalEntityPatterns,
  metaAlarmFilterPatternsToForm,
} from '@/helpers/entities/meta-alarm/rule/form';

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
