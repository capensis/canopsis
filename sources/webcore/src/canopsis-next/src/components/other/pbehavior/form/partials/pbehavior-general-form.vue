<template lang="pug">
  div
    v-layout(row)
      v-text-field(
        data-test="pbehaviorFormName",
        v-field="form.name",
        v-validate="'required'",
        :label="$t('modals.createPbehavior.fields.name')",
        :error-messages="errors.collect('name')",
        name="name"
      )
    v-layout(data-test="startDateTimePicker", row)
      date-time-picker-field(
        v-validate="tstartRules",
        v-field="form.tstart",
        v-validate="tstartRules",
        :label="$t('modals.createPbehavior.fields.start')",
        name="tstart"
      )
    v-layout(data-test="stopDateTimePicker", row)
      date-time-picker-field(
        v-validate="tstopRules",
        v-field="form.tstop",
        v-validate="tstopRules",
        :label="$t('modals.createPbehavior.fields.stop')",
        name="tstop"
      )
    v-layout(v-if="!noFilter", row)
      v-btn.primary(
        data-test="pbehaviorFilterButton",
        type="button",
        @click="showCreateFilterModal"
      ) {{ $t('common.filter') }}
    r-rule-form(:value="form.rrule", @input="updateField('rrule', $event)")
    v-layout(row)
      v-combobox(
        data-test="pbehaviorReason",
        v-field="form.reason",
        v-validate="'required'",
        :label="$t('modals.createPbehavior.fields.reason')",
        :items="reasons",
        :error-messages="errors.collect('reason')",
        name="reason"
      )
    v-layout(data-test="pbehaviorType", row)
      v-select(
        v-field="form.type_",
        v-validate="'required'",
        :label="$t('modals.createPbehavior.fields.type')",
        :items="types",
        :error-messages="errors.collect('type')",
        name="type"
      )
</template>

<script>
import moment from 'moment';
import { ENTITIES_TYPES, MODALS, PAUSE_REASONS, PBEHAVIOR_TYPES, DATETIME_FORMATS } from '@/constants';

import authMixin from '@/mixins/auth';
import formMixin from '@/mixins/form';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import RRuleForm from '@/components/forms/rrule.vue';

export default {
  inject: ['$validator'],
  components: {
    DateTimePickerField,
    RRuleForm,
  },
  mixins: [authMixin, formMixin],
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
  },
  computed: {
    reasons() {
      return Object.values(PAUSE_REASONS);
    },

    types() {
      return Object.values(PBEHAVIOR_TYPES);
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
  methods: {
    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          hiddenFields: ['title'],
          entitiesType: ENTITIES_TYPES.pbehavior,
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
