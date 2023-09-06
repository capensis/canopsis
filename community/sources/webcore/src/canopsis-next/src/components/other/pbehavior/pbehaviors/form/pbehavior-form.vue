<template lang="pug">
  v-tabs(slider-color="primary", centered, fixed-tabs)
    v-tab {{ $t('common.general') }}
    v-tab-item
      v-layout.py-3(column)
        pbehavior-general-form(
          v-field="form",
          :no-enabled="noEnabled",
          :with-start-on-trigger="withStartOnTrigger"
        )
        c-enabled-color-picker-field(
          v-field="form.color",
          :label="$t('modals.createPbehavior.steps.color.label')",
          row
        )
        c-collapse-panel.mb-2(:title="$t('recurrenceRule.title')")
          recurrence-rule-form(v-field="form.rrule", :start="form.tstart")
          pbehavior-recurrence-rule-exceptions-field.mt-2(
            v-field="form.exdates",
            :exceptions="form.exceptions",
            with-exdate-type,
            @update:exceptions="updateExceptions"
          )
        c-collapse-panel.mt-2(v-if="!noComments", :title="$tc('common.comment', 2)")
          pbehavior-comments-field(v-field="form.comments")
    v-tab {{ $tc('common.pattern', 2) }}
    v-tab-item
      v-layout.py-3(row, justify-center)
        v-flex(xs12)
          c-patterns-field(
            v-field="form.patterns",
            with-entity,
            some-required
          )
</template>

<script>
import { formMixin } from '@/mixins/form';

import RecurrenceRuleForm from '@/components/forms/recurrence-rule/recurrence-rule-form.vue';
import PbehaviorRecurrenceRuleExceptionsField from '@/components/other/pbehavior/exceptions/fields/pbehavior-recurrence-rule-exceptions-field.vue';

import PbehaviorCommentsField from '../fields/pbehavior-comments-field.vue';
import PbehaviorFilterField from '../fields/pbehavior-filter-field.vue';

import PbehaviorGeneralForm from './pbehavior-general-form.vue';

export default {
  inject: ['$validator'],
  components: {
    RecurrenceRuleForm,
    PbehaviorRecurrenceRuleExceptionsField,
    PbehaviorFilterField,
    PbehaviorGeneralForm,
    PbehaviorCommentsField,
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
  methods: {
    updateExceptions(exceptions) {
      this.updateField('exceptions', exceptions);
    },
  },
};
</script>
