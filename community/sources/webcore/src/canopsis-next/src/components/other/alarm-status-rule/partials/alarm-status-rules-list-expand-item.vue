<template>
  <v-tabs
    color="secondary lighten-1"
    slider-color="primary"
    dark="dark"
    centered="centered"
  >
    <v-tab>{{ $t('common.description') }}</v-tab>
    <v-tab-item>
      <v-layout
        class="py-3"
      >
        <v-textarea
          class="my-2 mx-4 pa-0"
          :value="rule.description"
          readonly="readonly"
          auto-grow="auto-grow"
          outlined
          hide-details="hide-details"
          dark="dark"
        />
      </v-layout>
    </v-tab-item>
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
              <alarm-status-rule-patterns-form
                :form="patterns"
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
import { OLD_PATTERNS_FIELDS, PATTERNS_FIELDS } from '@/constants';

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
  computed: {
    patterns() {
      return filterPatternsToForm(
        this.rule,
        [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity],
        [OLD_PATTERNS_FIELDS.alarm, OLD_PATTERNS_FIELDS.entity],
      );
    },
  },
};
</script>
