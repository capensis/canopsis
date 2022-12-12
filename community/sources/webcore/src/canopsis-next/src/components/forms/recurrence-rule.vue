<template lang="pug">
  div
    v-layout(row)
      v-tabs.recurrence-rule-tabs(slider-color="primary", fixed-tabs, centered)
        v-tab {{ $t('recurrenceRule.tabs.simple') }}
        v-tab-item
          v-layout(column)
            v-select(
              v-model="form.recurrenceRuleOptions.freq",
              :label="$t('recurrenceRule.freq')",
              :items="frequencies",
              @input="changeRecurrenceRuleOption"
            )
            v-chip-group(
              v-model="form.recurrenceRuleOptions.byweekday",
              :items="weekDays",
              :label="$t('recurrenceRule.byweekday')",
              multiple,
              @change="changeRecurrenceRuleOption"
            )
            c-number-field(
              v-model="form.recurrenceRuleOptions.count",
              :label="$t('recurrenceRule.count')",
              :min="1",
              name="count",
              @input="changeRecurrenceRuleOption"
            )
            c-number-field(
              v-model="form.recurrenceRuleOptions.interval",
              :label="$t('recurrenceRule.interval')",
              :min="1",
              name="interval",
              @input="changeRecurrenceRuleOption"
            )

        v-tab {{ $t('recurrenceRule.tabs.advanced') }}
        v-tab-item
          v-layout(column)
            v-select(
              v-model="form.recurrenceRuleOptions.wkst",
              :label="$t('recurrenceRule.wkst')",
              :items="weekDays",
              @input="changeRecurrenceRuleOption"
            )
            v-select(
              v-model="form.recurrenceRuleOptions.bymonth",
              :label="$t('recurrenceRule.bymonth')",
              :items="months",
              multiple,
              chips,
              @input="changeRecurrenceRuleOption"
            )
            v-tooltip(max-width="250", left)
              template(#activator="{ on }")
                v-text-field(
                  v-on="on",
                  v-model="form.advancedRecurrenceRuleOptions.bysetpos",
                  v-validate="{ regex: advancedFieldRegex }",
                  :label="$t('recurrenceRule.bysetpos')",
                  :error-messages="errors.collect('bysetpos')",
                  :hint="$t('recurrenceRule.advancedHint')",
                  name="bysetpos",
                  persistent-hint,
                  @input="changeRecurrenceRuleOption"
                )
              span {{ $t('recurrenceRule.tooltips.bysetpos') }}
            v-tooltip(max-width="250", left)
              template(#activator="{ on }")
                v-text-field(
                  v-on="on",
                  v-model="form.advancedRecurrenceRuleOptions.bymonthday",
                  v-validate="{ regex: advancedFieldRegex }",
                  :label="$t('recurrenceRule.bymonthday')",
                  :error-messages="errors.collect('bymonthday')",
                  :hint="$t('recurrenceRule.advancedHint')",
                  name="bymonthday",
                  persistent-hint,
                  @input="changeRecurrenceRuleOption"
                )
              span {{ $t('recurrenceRule.tooltips.bymonthday') }}
            v-tooltip(max-width="250", left)
              template(#activator="{ on }")
                v-text-field(
                  v-on="on",
                  v-model="form.advancedRecurrenceRuleOptions.byyearday",
                  v-validate="{ regex: advancedFieldRegex }",
                  :label="$t('recurrenceRule.byyearday')",
                  :error-messages="errors.collect('byyearday')",
                  :hint="$t('recurrenceRule.advancedHint')",
                  name="byyearday",
                  persistent-hint,
                  @input="changeRecurrenceRuleOption"
                )
              span {{ $t('recurrenceRule.tooltips.byyearday') }}
            v-tooltip(max-width="250", left)
              template(#activator="{ on }")
                v-text-field(
                  v-on="on",
                  v-model="form.advancedRecurrenceRuleOptions.byweekno",
                  v-validate="{ regex: advancedFieldRegex }",
                  :label="$t('recurrenceRule.byweekno')",
                  :error-messages="errors.collect('byweekno')",
                  :hint="$t('recurrenceRule.advancedHint')",
                  name="byweekno",
                  persistent-hint,
                  @input="changeRecurrenceRuleOption"
                )
              span {{ $t('recurrenceRule.tooltips.byweekno') }}
            v-tooltip(max-width="250", left)
              template(#activator="{ on }")
                v-text-field(
                  v-on="on",
                  v-model="form.advancedRecurrenceRuleOptions.byhour",
                  v-validate="{ regex: advancedFieldRegex }",
                  :label="$t('recurrenceRule.byhour')",
                  :error-messages="errors.collect('byhour')",
                  :hint="$t('recurrenceRule.advancedHint')",
                  name="byhour",
                  persistent-hint,
                  @input="changeRecurrenceRuleOption"
                )
              span {{ $t('recurrenceRule.tooltips.byhour') }}div
            v-tooltip(max-width="250", left)
              template(#activator="{ on }")
                v-text-field(
                  v-on="on",
                  v-model="form.advancedRecurrenceRuleOptions.byminute",
                  v-validate="{ regex: advancedFieldRegex }",
                  :label="$t('recurrenceRule.byminute')",
                  :error-messages="errors.collect('byminute')",
                  :hint="$t('recurrenceRule.advancedHint')",
                  name="byminute",
                  persistent-hint,
                  @input="changeRecurrenceRuleOption"
                )
              span {{ $t('recurrenceRule.tooltips.byminute') }}
            v-tooltip(max-width="250", left)
              template(#activator="{ on }")
                v-text-field(
                  v-on="on",
                  v-model="form.advancedRecurrenceRuleOptions.bysecond",
                  v-validate="{ regex: advancedFieldRegex }",
                  :label="$t('recurrenceRule.bysecond')",
                  :error-messages="errors.collect('bysecond')",
                  :hint="$t('recurrenceRule.advancedHint')",
                  name="bysecond",
                  persistent-hint,
                  @input="changeRecurrenceRuleOption"
                )
              span {{ $t('recurrenceRule.tooltips.bysecond') }}
    v-layout(row)
      v-flex(xs2)
        strong {{ $t('common.recurrence') }}
      v-flex(xs10)
        p {{ recurrenceRuleString }}
    v-layout(row)
      v-flex(xs2)
        strong {{ $t('common.summary') }}
      v-flex(xs10)
        p {{ recurrenceRuleText }}
    v-layout(row)
      v-alert(:value="errors.has('recurrenceRule')", type="error")
        span {{ errors.first('recurrenceRule') }}
</template>

<script>
import { RRule, rrulestr } from 'rrule';
import { mapValues, pickBy } from 'lodash';

import { recurrenceRuleToFormAdvancedOptions, recurrenceRuleToFormOptions } from '@/helpers/forms/recurrence-rule';

export default {
  inject: ['$validator'],
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
      recurrenceRule = new RRule({
        freq: RRule.DAILY,
      });
    }

    return {
      recurrenceRuleObject: recurrenceRule,
      form: {
        recurrenceRuleOptions: recurrenceRuleToFormOptions(recurrenceRule.origOptions),
        advancedRecurrenceRuleOptions: recurrenceRuleToFormAdvancedOptions(recurrenceRule.origOptions),
      },
    };
  },
  computed: {
    advancedFieldRegex() {
      return /^(-?(?:[1-9]|[1-9]\d|[12]\d\d|3[0-6][0-6]?))(,-?(?:[1-9]|[1-9]\d|[12]\d\d|3[0-6][0-6]?))*,?$/;
    },

    frequencies() {
      return [
        { text: this.$t('common.timeFrequencies.secondly'), value: RRule.SECONDLY },
        { text: this.$t('common.timeFrequencies.minutely'), value: RRule.MINUTELY },
        { text: this.$t('common.timeFrequencies.hourly'), value: RRule.HOURLY },
        { text: this.$t('common.timeFrequencies.daily'), value: RRule.DAILY },
        { text: this.$t('common.timeFrequencies.weekly'), value: RRule.WEEKLY },
        { text: this.$t('common.timeFrequencies.monthly'), value: RRule.MONTHLY },
        { text: this.$t('common.timeFrequencies.yearly'), value: RRule.YEARLY },
      ];
    },

    weekDays() {
      return [
        { text: this.$t('common.weekDays.monday'), value: RRule.MO.weekday },
        { text: this.$t('common.weekDays.tuesday'), value: RRule.TU.weekday },
        { text: this.$t('common.weekDays.wednesday'), value: RRule.WE.weekday },
        { text: this.$t('common.weekDays.thursday'), value: RRule.TH.weekday },
        { text: this.$t('common.weekDays.friday'), value: RRule.FR.weekday },
        { text: this.$t('common.weekDays.saturday'), value: RRule.SA.weekday },
        { text: this.$t('common.weekDays.sunday'), value: RRule.SU.weekday },
      ];
    },

    months() {
      return [
        { text: this.$t('common.months.january'), value: 1 },
        { text: this.$t('common.months.february'), value: 2 },
        { text: this.$t('common.months.march'), value: 3 },
        { text: this.$t('common.months.april'), value: 4 },
        { text: this.$t('common.months.may'), value: 5 },
        { text: this.$t('common.months.june'), value: 6 },
        { text: this.$t('common.months.july'), value: 7 },
        { text: this.$t('common.months.august'), value: 8 },
        { text: this.$t('common.months.september'), value: 9 },
        { text: this.$t('common.months.october'), value: 10 },
        { text: this.$t('common.months.november'), value: 11 },
        { text: this.$t('common.months.december'), value: 12 },
      ];
    },

    recurrenceRuleString() {
      return this.recurrenceRuleObject.toString();
    },

    recurrenceRuleText() {
      return this.recurrenceRuleObject.toText();
    },
  },
  mounted() {
    this.changeRecurrenceRuleOption();
  },
  methods: {
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

          this.$emit('input', this.recurrenceRuleString.replace(/.*RRULE:/, ''));
        }
      } catch (err) {
        console.error(err);
      }
    },
  },
};
</script>

<style scoped>
  .recurrence-rule-tabs {
    width: 100%;
  }

  p {
    -ms-word-break: break-all;
    word-break: break-all;
  }
</style>
