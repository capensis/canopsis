<template lang="pug">
  div
    v-layout(align-center)
      v-text-field(
        v-field="form.general.name",
        v-validate="'required'",
        :label="$t('modals.createPbehavior.steps.general.fields.name')",
        :error-messages="errors.collect('name')",
        name="name",
        data-test="pbehaviorFormName"
      )
      color-picker-field.ma-2(v-model="form.general.color")
    v-layout
      v-switch(
        v-field="form.general.enabled",
        :label="$t('modals.createPbehavior.steps.general.fields.enabled')",
        color="primary",
        hide-details
      )
    v-layout.mt-3(wrap)
      v-flex(xs12)
        v-layout(wrap, justify-space-between)
          v-flex(data-test="startDateTimePicker", xs4)
            date-time-picker-field(
              v-validate="tstartRules",
              :value="form.general.tstart",
              :label="$t('modals.createPbehavior.steps.general.fields.start')",
              name="tstart",
              @input="updateField('tstart', $event)"
            )
          v-flex(data-test="stopDateTimePicker", xs4)
            date-time-picker-field(
              v-validate="tstopRules",
              :value="form.general.tstop",
              :label="$t('modals.createPbehavior.steps.general.fields.stop')",
              name="tstop",
              @input="updateField('tstop', $event)"
            )
          v-flex(xs3)
            v-autocomplete(
              v-field="form.timezone",
              v-validate="'required'",
              :items="timezones",
              :label="$t('modals.createPbehavior.steps.general.fields.timezone')",
              name="timezone"
            )
    v-tabs(fixed-tabs, slider-color="primary")
      v-tab {{ $t('modals.createPbehavior.steps.general.title') }}
      v-tab {{ $t('modals.createPbehavior.steps.filter.title') }}
      v-tab {{ $t('modals.createPbehavior.steps.rrule.title') }}
      v-tab-item
        pbehavior-general-form(v-model="form.general")
        pbehavior-comments-form(v-model="form.comments")
</template>

<script>
import moment from 'moment-timezone';

import { DATETIME_FORMATS } from '@/constants';

import formValidationHeaderMixin from '@/mixins/form/validation-header';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import ColorPickerField from '@/components/forms/fields/color-picker.vue';

import PbehaviorGeneralForm from './pbehavior-general-form.vue';
import PbehaviorCommentsForm from './pbehavior-comments-form.vue';

export default {
  components: {
    ColorPickerField,
    PbehaviorCommentsForm,
    DateTimePickerField,
    PbehaviorGeneralForm,
  },
  mixins: [formValidationHeaderMixin],
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
      stepper: 1,
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

    timezones() {
      return moment.tz.names();
    },
  },
};
</script>
