<template>
  <div>
    <v-tabs
      v-model="activeTab"
      background-color="secondary lighten-1"
      slider-color="primary"
      dark
      centered
    >
      <v-tab>{{ $t('common.summary') }}</v-tab>
      <v-tab>{{ $tc('common.pattern', 2) }}</v-tab>
    </v-tabs>
    <v-layout class="pa-3 secondary lighten-2">
      <v-flex xs12>
        <v-card class="pa-3">
          <v-tabs-items
            class="pt-2"
            v-model="activeTab"
          >
            <v-tab-item>
              <v-flex
                offset-xs2
                xs8
              >
                <idle-rules-summary-tab :idle-rule="idleRule" />
              </v-flex>
            </v-tab-item>
            <v-tab-item>
              <v-flex
                offset-xs2
                xs8
              >
                <idle-rule-patterns-form
                  :form="patterns"
                  :is-entity-type="isEntityType"
                  readonly
                />
              </v-flex>
            </v-tab-item>
          </v-tabs-items>
        </v-card>
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
import { PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';
import { isIdleRuleEntityType } from '@/helpers/entities/idle-rule/form';

import IdleRulePatternsForm from '@/components/other/idle-rule/form/idle-rule-patterns-form.vue';

import IdleRulesSummaryTab from './idle-rules-summary-tab.vue';

export default {
  components: { IdleRulePatternsForm, IdleRulesSummaryTab },
  props: {
    idleRule: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      activeTab: 0,
    };
  },
  computed: {
    isEntityType() {
      return isIdleRuleEntityType(this.idleRule.type);
    },

    patterns() {
      return filterPatternsToForm(this.idleRule, [PATTERNS_FIELDS.entity, PATTERNS_FIELDS.alarm]);
    },
  },
};
</script>
