<template lang="pug">
  v-stepper(v-model="stepper", non-linear)
    v-stepper-header
      v-stepper-step.py-0(
        :complete="stepper > 1",
        step="1",
        editable,
        :rules="[() => !hasGeneralFormAnyError]"
      ) {{ $t('modals.createPbehavior.steps.general.title') }}
        small(v-if="hasGeneralFormAnyError") {{ $t('modals.createPbehavior.errors.invalid') }}
      template(v-if="!noFilter")
        v-divider
        v-stepper-step.py-0(
          :complete="stepper > 2",
          step="2",
          editable,
          :rules="[() => !hasFilterEditorAnyError]"
        ) {{ $t('modals.createPbehavior.steps.filter.title') }}
          small(v-if="hasFilterEditorAnyError") {{ $t('modals.createPbehavior.errors.invalid') }}
          small.font-italic.font-weight-light(v-else) {{ $t('common.optional') }}
      v-divider
      v-stepper-step.py-0(
        :complete="stepper > 3",
        step="3",
        editable
      ) {{ $t('modals.createPbehavior.steps.rrule.title') }}
        small.font-italic.font-weight-light {{ $t('common.optional') }}
      v-divider
      v-stepper-step.py-0(
        :complete="stepper > 4",
        step="4",
        editable
      ) {{ $t('modals.createPbehavior.steps.comments.title') }}
        small.font-italic.font-weight-light {{ $t('common.optional') }}
    v-stepper-items
      v-stepper-content(step="1")
        v-card
          v-card-text
            pbehavior-general-form(v-field="form.general", ref="pbehaviorGeneralForm")
      v-stepper-content(step="2")
        v-card
          v-card-text
            filter-editor(v-field="form.filter", required)
      v-stepper-content(step="3")
        v-card
          v-card-text
            r-rule-form(v-field="form.general.rrule")
            pbehavior-exdates-form(v-if="form.general.rrule", v-field="form.exdate")
      v-stepper-content(step="4")
        v-card
          v-card-text
            pbehavior-comments-form(v-field="form.comments")
</template>

<script>
import RRuleForm from '@/components/forms/rrule.vue';
import FilterEditor from '@/components/other/filter/editor/filter-editor.vue';

import PbehaviorGeneralForm from './partials/pbehavior-general-form.vue';
import PbehaviorCommentsForm from './partials/pbehavior-comments-form.vue';
import PbehaviorExdatesForm from './partials/pbehavior-exdates-form.vue';

export default {
  components: {
    RRuleForm,
    FilterEditor,
    PbehaviorGeneralForm,
    PbehaviorCommentsForm,
    PbehaviorExdatesForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  inject: ['$validator'],
  props: {
    form: {
      type: Object,
      required: true,
    },
    noFilter: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      stepper: 1,
      hasGeneralFormAnyError: false,
    };
  },
  computed: {
    hasFilterEditorAnyError() {
      return this.errors.has('filter');
    },
  },
  mounted() {
    this.$watch(() => this.$refs.pbehaviorGeneralForm.hasAnyError, (value) => {
      this.hasGeneralFormAnyError = value;
    });
  },
};
</script>
