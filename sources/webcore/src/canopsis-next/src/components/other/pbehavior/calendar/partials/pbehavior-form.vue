<template lang="pug">
  div
    v-layout(align-center)
      v-text-field(
        v-field="form.name",
        v-validate="'required'",
        :label="$t('modals.createPbehavior.steps.general.fields.name')",
        :error-messages="errors.collect('name')",
        name="name",
        data-test="pbehaviorFormName"
      )
    v-layout
      v-switch(
        v-field="form.enabled",
        :label="$t('modals.createPbehavior.steps.general.fields.enabled')",
        color="primary",
        hide-details
      )
    v-layout.mt-3(wrap)
      v-flex(xs12)
        v-layout(wrap, justify-space-between)
          v-flex(data-test="startDateTimePicker", xs6)
            date-time-picker-field(
              v-validate="tstartRules",
              :value="form.tstart",
              :label="$t('modals.createPbehavior.steps.general.fields.start')",
              name="tstart",
              @input="updateField('tstart', $event)"
            )
          v-flex(data-test="stopDateTimePicker", xs6)
            date-time-picker-field(
              v-validate="tstopRules",
              :value="form.tstop",
              :label="$t('modals.createPbehavior.steps.general.fields.stop')",
              name="tstop",
              @input="updateField('tstop', $event)"
            )
    v-tabs(v-model="activeTab", fixed-tabs, slider-color="primary")
      v-tab {{ $t('modals.createPbehavior.steps.general.title') }}
      v-tab {{ $t('modals.createPbehavior.steps.rrule.title') }}
      v-tab-item
        pbehavior-general-form(v-field="form")
      v-tab-item
        r-rule-form(v-field="form.rrule")
        pbehavior-exception-dates-form(v-if="form.rrule", v-field="form.exdate")
</template>

<script>
import moment from 'moment-timezone';

import { DATETIME_FORMATS } from '@/constants';

import formValidationHeaderMixin from '@/mixins/form/validation-header';
import formMixin from '@/mixins/form/object';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import FilterEditor from '@/components/other/filter/editor/filter-editor.vue';
import RRuleForm from '@/components/forms/rrule.vue';

import PbehaviorGeneralForm from './pbehavior-general-form.vue';
import PbehaviorCommentsForm from './pbehavior-comments-form.vue';
import PbehaviorExceptionDatesForm from './pbehavior-exception-dates-form.vue';

export default {
  components: {
    RRuleForm,
    FilterEditor,
    DateTimePickerField,
    PbehaviorGeneralForm,
    PbehaviorCommentsForm,
    PbehaviorExceptionDatesForm,
  },
  mixins: [formValidationHeaderMixin, formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  inject: ['$validator'],
  props: {
    form: {
      type: Object,
      default: () => ({
        general: {},
      }),
    },
    noFilter: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      activeTab: 0,
      hasGeneralFormAnyError: false,
    };
  },
  computed: {
    hasFilterEditorAnyError() {
      return this.errors.has('filter');
    },

    tstartRules() {
      return {
        required: true,
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    tstopRules() {
      const rules = { required: true };

      if (this.form.tstart) {
        rules.after = [moment(this.form.tstart).format(DATETIME_FORMATS.dateTimePicker)];
        rules.date_format = DATETIME_FORMATS.veeValidateDateTimeFormat;
      }

      return rules;
    },
  },
};
</script>
