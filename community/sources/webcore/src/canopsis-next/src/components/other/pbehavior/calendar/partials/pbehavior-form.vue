<template lang="pug">
  v-layout(column)
    pbehavior-general-form(
      v-field="form",
      :no-enabled="noEnabled",
      :with-start-on-trigger="withStartOnTrigger"
    )
    pbehavior-comments-form(v-if="!noComments", v-field="form.comments")
    pbehavior-filter-field(v-if="!noPattern", v-field="form.patterns")
    pbehavior-recurrence-rule-field(v-field="form", with-exdate-type)
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
    noPattern: {
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
