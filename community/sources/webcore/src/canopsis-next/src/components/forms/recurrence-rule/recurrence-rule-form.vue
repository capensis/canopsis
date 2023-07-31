<template lang="pug">
  div.recurrence-rule-form
    v-tabs(v-model="activeTab", slider-color="primary", fixed-tabs, centered)
      v-tab {{ $t('recurrenceRule.tabs.simple') }}
      v-tab(:disabled="!isFrequencyEnabled") {{ $t('recurrenceRule.tabs.advanced') }}

    v-layout(column)
      v-layout(row)
        v-flex.pr-2(xs6)
          recurrence-rule-frequency-field(:value="form.recurrenceRuleOptions.freq", @input="updateFrequency")
          recurrence-rule-interval-field(v-if="isFrequencyEnabled", v-model="form.recurrenceRuleOptions")
        v-flex.pl-2(xs6)
          recurrence-rule-end-field(v-if="isFrequencyEnabled", v-model="form.recurrenceRuleOptions")

      recurrence-rule-weekday-field(
        v-if="!isAdvancedTab && isWeeklyFrequency",
        v-model="form.recurrenceRuleOptions.byweekday",
        chips
      )

    v-tabs-items(v-model="activeTab")
      v-tab-item
      v-tab-item(:disabled="!isFrequencyEnabled")
        recurrence-rule-weekday-field(v-model="form.recurrenceRuleOptions.wkst")
        recurrence-rule-weekday-field(v-if="!isYearlyFrequency", v-model="form.recurrenceRuleOptions.byweekday", chips)
        recurrence-rule-month-field(v-model="form.recurrenceRuleOptions.bymonth")
        v-layout(row, wrap)
          v-flex(v-for="(field, index) in advancedFields", :key="field", :class="`${index % 2 ? 'pl' : 'pr'}-2`", xs6)
            recurrence-rule-regex-field(
              v-model="form.advancedRecurrenceRuleOptions[field]",
              :label="$t(`recurrenceRule.${field}`)",
              :help-text="$t(`recurrenceRule.tooltips.${field}`)",
              :name="field"
            )
    template(v-if="isFrequencyEnabled")
      recurrence-rule-information(:rrule="recurrenceRuleString")
    c-alert(:value="errors.has('recurrenceRule')", type="error") {{ errors.first('recurrenceRule') }}
</template>

<script>
import { RRule, rrulestr } from 'rrule';
import { isNull, mapValues, pickBy } from 'lodash';

import { recurrenceRuleToFormAdvancedOptions, recurrenceRuleToFormOptions } from '@/helpers/entities/shared/recurrence-rule/form';

import RecurrenceRuleInformation from '@/components/common/reccurence-rule/recurrence-rule-information.vue';
import RecurrenceRuleRegexField from '@/components/forms/recurrence-rule/fields/recurrence-rule-regex-field.vue';
import RecurrenceRuleEndField from '@/components/forms/recurrence-rule/fields/recurrence-rule-end-field.vue';
import RecurrenceRuleIntervalField from '@/components/forms/recurrence-rule/fields/recurrence-rule-interval-field.vue';

import RecurrenceRuleMonthField from './fields/recurrence-rule-month-field.vue';
import RecurrenceRuleWeekdayField from './fields/recurrence-rule-weekday-field.vue';
import RecurrenceRuleFrequencyField from './fields/recurrence-rule-frequency-field.vue';

export default {
  inject: ['$validator'],
  components: {
    RecurrenceRuleIntervalField,
    RecurrenceRuleEndField,
    RecurrenceRuleRegexField,
    RecurrenceRuleMonthField,
    RecurrenceRuleWeekdayField,
    RecurrenceRuleFrequencyField,
    RecurrenceRuleInformation,
  },
  model: {
    prop: 'rrule',
    event: 'input',
  },
  props: {
    rrule: {
      type: String,
    },
  },
  data() {
    let recurrenceRule;

    if (this.rrule) {
      try {
        recurrenceRule = rrulestr(this.rrule);
      } catch (err) {
        console.error(err);
      }
    }

    if (!recurrenceRule) {
      recurrenceRule = new RRule();
    }

    return {
      activeTab: 0,
      recurrenceRuleObject: recurrenceRule,
      form: {
        recurrenceRuleOptions: recurrenceRuleToFormOptions(recurrenceRule.origOptions),
        advancedRecurrenceRuleOptions: recurrenceRuleToFormAdvancedOptions(recurrenceRule.origOptions),
      },
    };
  },
  computed: {
    isFrequencyEnabled() {
      return !isNull(this.form.recurrenceRuleOptions.freq);
    },

    isHourlyFrequency() {
      return this.form.recurrenceRuleOptions.freq === RRule.HOURLY;
    },

    isWeeklyFrequency() {
      return this.form.recurrenceRuleOptions.freq === RRule.WEEKLY;
    },

    isMonthlyFrequency() {
      return this.form.recurrenceRuleOptions.freq === RRule.MONTHLY;
    },

    isYearlyFrequency() {
      return this.form.recurrenceRuleOptions.freq === RRule.YEARLY;
    },

    isAdvancedTab() {
      return this.activeTab === 1;
    },

    advancedFields() {
      const fields = ['bysetpos'];

      if (!this.isMonthlyFrequency) {
        fields.push('byyearday');
      }

      if (!this.isYearlyFrequency) {
        fields.push('bymonthday');
      }

      if (!this.isMonthlyFrequency && !this.isYearlyFrequency) {
        fields.push('byweekno');
      }

      if (this.isHourlyFrequency) {
        fields.push('byhour');
      }

      return fields;
    },

    recurrenceRuleString() {
      return this.recurrenceRuleObject.toString();
    },
  },
  watch: {
    form: {
      deep: true,
      handler() {
        this.changeRecurrenceRuleOption();
      },
    },
  },
  mounted() {
    this.changeRecurrenceRuleOption();
  },
  methods: {
    updateFrequency(frequency) {
      this.form.recurrenceRuleOptions.freq = frequency;

      if (this.isAdvancedTab && isNull(frequency)) {
        this.activeTab = 0;
      }
    },
    /**
     * For each changes in the form we call this function.
     * If RRule isn't valid then add error message to visible RRule field
     * Else remove errors and $emit changes
     */
    changeRecurrenceRuleOption() {
      try {
        this.recurrenceRuleObject = new RRule({
          ...pickBy(this.form.recurrenceRuleOptions, v => v !== ''),
          ...mapValues(this.form.advancedRecurrenceRuleOptions, o => o.split(',').filter(v => v)),
        });

        if (!this.errors.has('recurrenceRule') && !this.recurrenceRuleObject.isFullyConvertibleToText()) {
          this.errors.add({
            field: 'recurrenceRule',
            msg: this.$t('recurrenceRule.errors.main'),
          });
        } else {
          this.errors.remove('recurrenceRule');

          /** TODO: Should be used updateModel */
          this.$emit('input', this.recurrenceRuleString.replace(/.*RRULE:/, ''));
        }
      } catch (err) {
        this.$emit('input', '');
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.recurrence-rule-form {
  p {
    -ms-word-break: break-all;
    word-break: break-all;
  }
}
</style>
