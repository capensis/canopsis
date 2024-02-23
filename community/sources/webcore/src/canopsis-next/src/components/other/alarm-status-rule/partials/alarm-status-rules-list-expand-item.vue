<template>
  <v-tabs
    v-model="activeTab"
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $t('common.description') }}</v-tab>
    <v-tab>{{ $tc('common.pattern', 2) }}</v-tab>

    <v-tabs-items
      v-model="activeTab"
      mandatory
    >
      <v-tab-item>
        <v-layout class="py-3">
          <v-textarea
            :value="rule.description"
            class="my-2 mx-4 pa-0"
            readonly
            auto-grow
            outlined
            hide-details
            dark
          />
        </v-layout>
      </v-tab-item>
      <v-tab-item>
        <v-layout class="py-3">
          <v-flex
            xs12
            md8
            offset-md2
          >
            <v-card>
              <v-card-text>
                <alarm-status-rule-patterns-form
                  :form="patterns"
                  readonly
                />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
    </v-tabs-items>
  </v-tabs>
</template>

<script>
import { PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';

import AlarmStatusRulePatternsForm from '../form/alarm-status-rule-patterns-form.vue';

export default {
  components: {
    AlarmStatusRulePatternsForm,
  },
  props: {
    rule: {
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
    patterns() {
      return filterPatternsToForm(this.rule, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]);
    },
  },
};
</script>
