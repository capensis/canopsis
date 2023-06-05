<template lang="pug">
  v-tabs(color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $tc('common.pattern', 2) }}
    v-tab-item
      v-layout.py-3(row)
        v-flex(xs12, sm10, offset-sm1)
          v-card
            v-card-text
              pbehavior-patterns-form(:form="patterns", readonly)

    v-tab {{ $tc('common.entity', 2) }}
    v-tab-item(lazy)
      v-layout.py-3(row)
        v-flex(xs12, sm10, offset-sm1)
          v-card
            v-card-text
              pbehavior-entities(:pbehavior="pbehavior")

    v-tab {{ $tc('common.comment', 2) }}
    v-tab-item(lazy)
      v-layout.py-3(row)
        v-flex(xs12, sm10, offset-sm1)
          v-card
            v-card-text.pa-0
              pbehavior-comments(:comments="pbehavior.comments")

    template(v-if="pbehavior.rrule")
      v-tab {{ $t('common.recurrence') }}
      v-tab-item(lazy)
        v-layout.py-3(row)
          v-flex(xs12, sm10, offset-sm1)
            v-card
              v-card-text
                pbehavior-recurrence-rule(:pbehavior="pbehavior")
</template>

<script>
import { OLD_PATTERNS_FIELDS, PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/forms/filter';

import PbehaviorPatternsForm from '../form/pbehavior-patterns-form.vue';

import PbehaviorComments from './pbehavior-comments.vue';
import PbehaviorRecurrenceRule from './pbehavior-recurrence-rule.vue';
import PbehaviorEntities from './pbehavior-entities.vue';

export default {
  components: {
    PbehaviorPatternsForm,
    PbehaviorComments,
    PbehaviorRecurrenceRule,
    PbehaviorEntities,
  },
  props: {
    pbehavior: {
      type: Object,
      required: true,
    },
  },
  computed: {
    patterns() {
      return filterPatternsToForm(
        this.pbehavior,
        [PATTERNS_FIELDS.entity],
        [OLD_PATTERNS_FIELDS.mongoQuery],
      );
    },
  },
};
</script>
