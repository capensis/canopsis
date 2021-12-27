<template lang="pug">
  div
    v-layout(wrap)
      v-flex(xs12)
        pbehavior-general-form(
          v-field="form",
          :no-enabled="noEnabled",
          :with-start-on-trigger="withStartOnTrigger"
        )
      v-flex(v-if="!noComments", xs12)
        pbehavior-comments-form(v-field="form.comments")
      v-flex(v-if="!noFilter", xs12)
        pbehavior-filter-field(v-field="form")
      v-flex(xs12)
        pbehavior-recurrence-rule-field(v-field="form")
</template>

<script>
import { formMixin } from '@/mixins/form';

import PbehaviorGeneralForm from './pbehavior-general-form.vue';
import PbehaviorCommentsForm from './pbehavior-comments-form.vue';
import PbehaviorFilterField from './pbehavior-filter-field.vue';
import PbehaviorRecurrenceRuleField from './pbehavior-recurrence-rule-field.vue';

export default {
  inject: ['$validator'],
  components: {
    PbehaviorRecurrenceRuleField,
    PbehaviorFilterField,
    PbehaviorGeneralForm,
    PbehaviorCommentsForm,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    noFilter: {
      type: Boolean,
      default: false,
    },
    noEnabled: {
      type: Boolean,
      default: false,
    },
    noComments: {
      type: Boolean,
      default: false,
    },
    withStartOnTrigger: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
