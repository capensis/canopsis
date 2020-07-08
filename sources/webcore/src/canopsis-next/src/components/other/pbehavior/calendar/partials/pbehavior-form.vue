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
      pbehavior-general-form(v-field="form")
      v-flex(xs12)
        v-btn.ml-0.btn-filter(
          :color="errors.has('filter') ? 'error' : 'primary'",
          @click="showCreateFilterModal"
        ) {{ hasFilter ? 'Edit filter' : 'Add filter' }}
        v-tooltip(v-show="hasFilter", fixed, top)
          v-btn(slot="activator", icon)
            v-icon(color="grey darken-1") info
          span.pre {{ form.filter.filter | json }}
        v-alert(:value="errors.has('filter')", type="error") {{ errors.first('filter') }}
      v-flex(xs12)
        v-btn.ml-0(
          color="primary",
          @click="showCreateRRuleModal"
        )  {{ hasRRule ? 'Edit RRule' : 'Add RRule' }}
        v-tooltip(v-show="hasRRule", fixed, top)
          v-btn(slot="activator", icon)
            v-icon(color="grey darken-1") info
          span {{ form.rrule }}
</template>

<script>
import moment from 'moment-timezone';
import { isEmpty } from 'lodash';

import { DATETIME_FORMATS, MODALS } from '@/constants';

import formValidationHeaderMixin from '@/mixins/form/validation-header';
import formMixin from '@/mixins/form/object';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import FilterEditor from '@/components/other/filter/editor/filter-editor.vue';
import RRuleForm from '@/components/forms/rrule.vue';

import PbehaviorGeneralForm from './pbehavior-general-form.vue';
import PbehaviorCommentsForm from './pbehavior-comments-form.vue';
import PbehaviorExdatesForm from './pbehavior-exdates-form.vue';

export default {
  components: {
    RRuleForm,
    FilterEditor,
    DateTimePickerField,
    PbehaviorGeneralForm,
    PbehaviorCommentsForm,
    PbehaviorExdatesForm,
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
  computed: {
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
    hasFilter() {
      return this.form.filter && !isEmpty(this.form.filter.filter);
    },
    hasRRule() {
      return !isEmpty(this.form.rrule);
    },
  },
  created() {
    this.$validator.attach({
      name: 'filter',
      rules: 'required:true',
      getter: () => !isEmpty(this.form.filter),
      context: () => this,
      vm: this,
    });
  },
  methods: {
    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        dialogProps: {
          zIndex: 300,
        },
        config: {
          filter: this.form.filter,
          hiddenFields: ['title'],
          action: (filter) => {
            this.$emit('input', { ...this.form, filter });
            this.$nextTick(() => this.$validator.validate('filter'));
          },
        },
      });
    },
    showCreateRRuleModal() {
      this.$modals.show({
        name: MODALS.createRRule,
        dialogProps: {
          zIndex: 300,
        },
        config: {
          rrule: this.form.rrule,
          action: rrule => this.$emit('input', { ...this.form, rrule }),
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .pre {
    white-space: pre;
  }

  .btn-filter.error {
    animation: shake .6s cubic-bezier(.25,.8,.5,1);
  }

  @keyframes shake {
    10%, 90% {
      transform: translate3d(-1px, 0, 0);
    }

    20%, 80% {
      transform: translate3d(2px, 0, 0);
    }

    30%, 50%, 70% {
      transform: translate3d(-4px, 0, 0);
    }

    40%, 60% {
      transform: translate3d(4px, 0, 0);
    }
  }
</style>
