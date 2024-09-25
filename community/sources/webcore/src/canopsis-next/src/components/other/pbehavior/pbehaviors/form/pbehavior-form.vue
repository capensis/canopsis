<template lang="pug">
  v-layout(column)
    pbehavior-general-form(
      v-field="form",
      :no-enabled="noEnabled",
      :with-start-on-trigger="withStartOnTrigger",
      :name-label="nameLabel",
      :name-tooltip="nameTooltip"
    )
    pbehavior-comments-field(v-if="!noComments", v-field="form.comments")
    pbehavior-filter-field(v-if="!noPattern", v-field="form.patterns")
    pbehavior-recurrence-rule-field(v-field="form", with-exdate-type)
    c-enabled-color-picker-field(v-field="form.color", :label="$t('modals.createPbehavior.steps.color.label')")
</template>

<script>
import { formMixin } from '@/mixins/form';

import PbehaviorCommentsField from '../fields/pbehavior-comments-field.vue';
import PbehaviorFilterField from '../fields/pbehavior-filter-field.vue';
import PbehaviorRecurrenceRuleField from '../fields/pbehavior-recurrence-rule-field.vue';

import PbehaviorGeneralForm from './pbehavior-general-form.vue';

export default {
  inject: ['$validator'],
  components: {
    PbehaviorRecurrenceRuleField,
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
    nameLabel: {
      type: String,
      required: false,
    },
    nameTooltip: {
      type: String,
      required: false,
    },
  },
};
</script>
