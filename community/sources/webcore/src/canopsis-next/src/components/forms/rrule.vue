<template lang="pug">
  div
    v-layout(row)
      v-switch(
        v-model="showRRule",
        :label="$t('modals.createPbehavior.steps.rrule.fields.rRuleQuestion')",
        data-test="pbehaviorRuleSwitcher",
        color="primary"
      )
    template(v-if="showRRule")
      v-layout(row)
        v-tabs.r-rule-tabs(v-model="activeTab", centered, fixed-tabs, slider-color="primary")
          v-tab(data-test="pbehaviorSimple", href="#simple") {{ $t('rRule.tabs.simple') }}
          v-tab(data-test="pbehaviorAdvanced", href="#advanced") {{ $t('rRule.tabs.advanced') }}
          v-tab-item(value="simple")
            div
              div(data-test="pbehaviorFrequency")
                v-select(
                  v-model="form.rRuleOptions.freq",
                  :label="$t('rRule.fields.freq')",
                  :items="frequencies",
                  @input="changeRRuleOption"
                )
              div(data-test="pbehaviorByWeekDay")
                v-chip-group(
                  v-model="form.rRuleOptions.byweekday",
                  :items="weekDays",
                  :label="$t('rRule.fields.byweekday')",
                  multiple,
                  @change="changeRRuleOption"
                )
              div
                v-text-field(
                  v-model.number="form.rRuleOptions.count",
                  v-validate="'numeric|min_value:1'",
                  :label="$t('rRule.fields.count')",
                  :error-messages="errors.collect('count')",
                  data-test="pbehaviorRepeat",
                  type="number",
                  name="count",
                  min="1",
                  @input="changeRRuleOption"
                )
              div
                v-text-field(
                  v-model.number="form.rRuleOptions.interval",
                  v-validate="'numeric|min_value:1'",
                  :label="$t('rRule.fields.interval')",
                  :error-messages="errors.collect('interval')",
                  data-test="pbehaviorInterval",
                  type="number",
                  name="interval",
                  min="1",
                  @input="changeRRuleOption"
                )
          v-tab-item(value="advanced")
            div
              div(data-test="pbehaviorWeekStart")
                v-select(
                  v-model="form.rRuleOptions.wkst",
                  :label="$t('rRule.fields.wkst')",
                  :items="weekDays",
                  @input="changeRRuleOption"
                )
              div(data-test="pbehaviorByMonth")
                v-select(
                  v-model="form.rRuleOptions.bymonth",
                  :label="$t('rRule.fields.bymonth')",
                  :items="months",
                  multiple,
                  chips,
                  @input="changeRRuleOption"
                )
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      v-model="form.advancedRRuleOptions.bysetpos",
                      v-validate="{ regex: advancedFieldRegex }",
                      :label="$t('rRule.fields.bysetpos.label')",
                      :error-messages="errors.collect('bysetpos')",
                      :hint="$t('rRule.advancedHint')",
                      data-test="pbehaviorBySetPos",
                      name="bysetpos",
                      persistent-hint,
                      @input="changeRRuleOption"
                    )
                  span {{ $t('rRule.fields.bysetpos.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      v-model="form.advancedRRuleOptions.bymonthday",
                      v-validate="{ regex: advancedFieldRegex }",
                      :label="$t('rRule.fields.bymonthday.label')",
                      :error-messages="errors.collect('bymonthday')",
                      :hint="$t('rRule.advancedHint')",
                      data-test="pbehaviorByMonthDay",
                      name="bymonthday",
                      persistent-hint,
                      @input="changeRRuleOption"
                    )
                  span {{ $t('rRule.fields.bymonthday.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      v-model="form.advancedRRuleOptions.byyearday",
                      v-validate="{ regex: advancedFieldRegex }",
                      :label="$t('rRule.fields.byyearday.label')",
                      :error-messages="errors.collect('byyearday')",
                      :hint="$t('rRule.advancedHint')",
                      data-test="pbehaviorByYearDay",
                      name="byyearday",
                      persistent-hint,
                      @input="changeRRuleOption"
                    )
                  span {{ $t('rRule.fields.byyearday.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      v-model="form.advancedRRuleOptions.byweekno",
                      v-validate="{ regex: advancedFieldRegex }",
                      :label="$t('rRule.fields.byweekno.label')",
                      :error-messages="errors.collect('byweekno')",
                      :hint="$t('rRule.advancedHint')",
                      data-test="pbehaviorByWeekNo",
                      name="byweekno",
                      persistent-hint,
                      @input="changeRRuleOption"
                    )
                  span {{ $t('rRule.fields.byweekno.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      v-model="form.advancedRRuleOptions.byhour",
                      v-validate="{ regex: advancedFieldRegex }",
                      :label="$t('rRule.fields.byhour.label')",
                      :error-messages="errors.collect('byhour')",
                      :hint="$t('rRule.advancedHint')",
                      data-test="pbehaviorByHour",
                      name="byhour",
                      persistent-hint,
                      @input="changeRRuleOption"
                    )
                  span {{ $t('rRule.fields.byhour.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      v-model="form.advancedRRuleOptions.byminute",
                      v-validate="{ regex: advancedFieldRegex }",
                      :label="$t('rRule.fields.byminute.label')",
                      :error-messages="errors.collect('byminute')",
                      :hint="$t('rRule.advancedHint')",
                      data-test="pbehaviorByMinute",
                      name="byminute",
                      persistent-hint,
                      @input="changeRRuleOption"
                    )
                  span {{ $t('rRule.fields.byminute.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      v-model="form.advancedRRuleOptions.bysecond",
                      v-validate="{ regex: advancedFieldRegex }",
                      :label="$t('rRule.fields.bysecond.label')",
                      :error-messages="errors.collect('bysecond')",
                      :hint="$t('rRule.advancedHint')",
                      data-test="pbehaviorBySecond",
                      name="bysecond",
                      persistent-hint,
                      @input="changeRRuleOption"
                    )
                  span {{ $t('rRule.fields.bysecond.tooltip') }}
      v-layout(row)
        v-flex(xs2)
          strong {{ $t('rRule.stringLabel') }}
        v-flex(xs10)
          p {{ rRuleString }}
      v-layout(row)
        v-flex(xs2)
          strong {{ $t('rRule.textLabel') }}
        v-flex(xs10)
          p {{ rRuleText }}
      v-layout(row)
        v-alert(:value="errors.has('rRule')", type="error")
          span {{ errors.first('rRule') }}
</template>

<script>
import { RRule, rrulestr } from 'rrule';
import { mapValues, pickBy } from 'lodash';

import { recurrenceRuleToFormAdvancedOptions, recurrenceRuleToFormOptions } from '@/helpers/forms/rrule';

/**
 * RRule form component
 *
 * @event rrule#input
 * @type {Object|null} - RRule object
 */
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
    let rRule;

    if (this.rrule) {
      try {
        rRule = rrulestr(this.rrule);
      } catch (err) {
        console.warn(err);
      }
    }

    if (!rRule) {
      rRule = new RRule({
        freq: RRule.DAILY,
      });
    }

    return {
      activeTab: 'simple',
      showRRule: !!this.rrule,
      rRuleObject: rRule,
      form: {
        rRuleOptions: recurrenceRuleToFormOptions(rRule.origOptions),
        advancedRRuleOptions: recurrenceRuleToFormAdvancedOptions(rRule.origOptions),
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

    rRuleString() {
      return this.rRuleObject.toString();
    },

    rRuleText() {
      return this.rRuleObject.toText();
    },
  },
  watch: {
    showRRule(value) {
      if (!value) {
        this.errors.remove('rRule');
        this.$emit('input', '');
      } else {
        this.changeRRuleOption();
      }
    },
  },
  methods: {
    /**
     * For each changes in the form we call this function.
     * If RRule isn't valid then add error message to visible RRule field
     * Else remove errors and $emit changes
     */
    changeRRuleOption() {
      try {
        this.rRuleObject = new RRule({
          ...pickBy(this.form.rRuleOptions, v => v !== ''),
          ...mapValues(this.form.advancedRRuleOptions, o => o.split(',').filter(v => v)),
        });

        if (!this.errors.has('rRule') && !this.rRuleObject.isFullyConvertibleToText()) {
          this.errors.add({
            field: 'rRule',
            msg: this.$t('rRule.errors.main'),
          });
        } else {
          this.errors.remove('rRule');

          if (this.showRRule) {
            this.$emit('input', this.rRuleString.replace(/.*RRULE:/, ''));
          }
        }
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
</script>

<style scoped>
  .r-rule-tabs {
    width: 100%;
  }

  p {
    -ms-word-break: break-all;
    word-break: break-all;
  }
</style>
