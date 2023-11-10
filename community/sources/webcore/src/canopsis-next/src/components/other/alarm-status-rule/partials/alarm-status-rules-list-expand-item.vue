<template lang="pug">
  v-tabs(color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $t('common.description') }}
    v-tab-item
      v-layout.py-3(row)
        v-textarea.my-2.mx-4.pa-0(
          :value="rule.description",
          readonly,
          auto-grow,
          outline,
          hide-details,
          dark
        )
    v-tab {{ $tc('common.pattern', 2) }}
    v-tab-item
      v-layout.py-3
        v-flex(xs12, md8, offset-md2)
          v-card
            v-card-text
              alarm-status-rule-patterns-form(:form="patterns", readonly)
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
  computed: {
    patterns() {
      return filterPatternsToForm(this.rule, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]);
    },
  },
};
</script>
