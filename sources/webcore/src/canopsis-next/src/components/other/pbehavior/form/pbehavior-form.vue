<template lang="pug">
  div
    v-layout(row)
      v-text-field(
      v-validate="'required'",
      :value="form.name",
      :label="$t('modals.createPbehavior.fields.name')",
      :error-messages="errors.collect('name')",
      name="name",
      @input="updateField('name', $event)"
      )
    v-layout(row)
      date-time-picker-field(
      v-validate="'required'",
      :value="form.tstart",
      :label="$t('modals.createPbehavior.fields.start')",
      name="tstart",
      @input="updateField('tstart', $event)"
      )
    v-layout(row)
      date-time-picker-field(
      v-validate="tstopRules",
      :value="form.tstop",
      :label="$t('modals.createPbehavior.fields.stop')",
      name="tstop",
      @input="updateField('tstop', $event)"
      )
    v-layout(row)
      v-btn.primary(type="button", @click="showCreateFilterModal") {{ $t('common.filter') }}
    v-layout(row)
      r-rule-form(:value="form.rrule", @input="updateField('rrule', $event)")
    v-layout(row)
      pbehavior-exdate-form(v-model="form.exdate")
    v-layout(row)
      v-combobox(
      v-validate="'required'",
      :value="form.reason",
      :label="$t('modals.createPbehavior.fields.reason')",
      :items="reasons",
      :error-messages="errors.collect('reason')",
      name="reason",
      @input="updateField('reason', $event)"
      )
    v-layout(row)
      v-select(
      v-validate="'required'",
      :value="form.type_",
      :label="$t('modals.createPbehavior.fields.type')",
      :items="types",
      :error-messages="errors.collect('type')",
      name="type",
      @input="updateField('type_', $event)"
      )
</template>

<script>
import moment from 'moment';
import { MODALS, PAUSE_REASONS, PBEHAVIOR_TYPES, DATETIME_FORMATS } from '@/constants';

import authMixin from '@/mixins/auth';
import formMixin from '@/mixins/form';
import modalMixin from '@/mixins/modal';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import RRuleForm from '@/components/forms/rrule.vue';
import PbehaviorExdateForm from './pbehavior-exdate-form.vue';

export default {
  inject: ['$validator'],
  components: {
    DateTimePickerField,
    RRuleForm,
    PbehaviorExdateForm,
  },
  mixins: [authMixin, formMixin, modalMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  computed: {
    reasons() {
      return Object.values(PAUSE_REASONS);
    },

    types() {
      return Object.values(PBEHAVIOR_TYPES);
    },

    tstopRules() {
      const rules = { required: true };

      if (this.form.tstart) {
        rules.after = [moment(this.form.tstart).format(DATETIME_FORMATS.dateTimePicker)];
        rules.date_format = DATETIME_FORMATS.dateTimePicker;
      }

      return rules;
    },
  },
  methods: {
    showCreateFilterModal() {
      this.showModal({
        name: MODALS.createFilter,
        config: {
          hiddenFields: ['title'],
          filter: {
            filter: this.form.filter || {},
          },
          action: ({ filter }) => this.updateField('filter', filter),
        },
      });
    },
  },
};
</script>
