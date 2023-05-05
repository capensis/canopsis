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
import { OLD_PATTERNS_FIELDS, PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/forms/filter';

import AlarmStatusRulePatternsForm from '../form/partials/alarm-status-rule-patterns-form.vue';

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
