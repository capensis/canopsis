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
import { mapValues, pickBy, isArray } from 'lodash';

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
        rRuleOptions: this.rRuleToFormOptions(rRule),
        advancedRRuleOptions: this.rRuleToFormAdvancedOptions(rRule),
      },
    };
  },
  computed: {
    advancedFieldRegex() {
      return /^(-?(?:[1-9]|[1-9]\d|[12]\d\d|3[0-6][0-6]?))(,-?(?:[1-9]|[1-9]\d|[12]\d\d|3[0-6][0-6]?))*,?$/;
    },

    frequencies() {
      return [
        { text: 'Secondly', value: RRule.SECONDLY },
        { text: 'Minutely', value: RRule.MINUTELY },
        { text: 'Hourly', value: RRule.HOURLY },
        { text: 'Daily', value: RRule.DAILY },
        { text: 'Weekly', value: RRule.WEEKLY },
        { text: 'Monthly', value: RRule.MONTHLY },
        { text: 'Yearly', value: RRule.YEARLY },
      ];
    },

    weekDays() {
      return [
        { text: 'Monday', value: RRule.MO.weekday },
        { text: 'Tuesday', value: RRule.TU.weekday },
        { text: 'Wednesday', value: RRule.WE.weekday },
        { text: 'Thursday', value: RRule.TH.weekday },
        { text: 'Friday', value: RRule.FR.weekday },
        { text: 'Saturday', value: RRule.SA.weekday },
        { text: 'Sunday', value: RRule.SU.weekday },
      ];
    },

    months() {
      return [
        { text: 'January', value: 1 },
        { text: 'February', value: 2 },
        { text: 'March', value: 3 },
        { text: 'April', value: 4 },
        { text: 'May', value: 5 },
        { text: 'June', value: 6 },
        { text: 'July', value: 7 },
        { text: 'August', value: 8 },
        { text: 'September', value: 9 },
        { text: 'October', value: 10 },
        { text: 'November', value: 11 },
        { text: 'December', value: 12 },
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
    rRuleToFormOptions(rRule) {
      const { origOptions } = rRule;

      return {
        freq: origOptions.freq || '',
        count: origOptions.count || '',
        interval: origOptions.interval || '',
        byweekday: origOptions.byweekday ? origOptions.byweekday.map(v => v.weekday) : [],
        wkst: origOptions.wkst ? origOptions.wkst.weekday : '',
        bymonth: origOptions.bymonth || [],
      };
    },

    rRuleToFormAdvancedOptions(rRule) {
      const { origOptions } = rRule;
      const optionPreparer = (v) => {
        if (v) {
          return (isArray(v) ? v.join(',') : String(v));
        }

        return '';
      };

      return {
        bysetpos: optionPreparer(origOptions.bysetpos),
        bymonthday: optionPreparer(origOptions.bymonthday),
        byyearday: optionPreparer(origOptions.byyearday),
        byweekno: optionPreparer(origOptions.byweekno),
        byhour: optionPreparer(origOptions.byhour),
        byminute: optionPreparer(origOptions.byminute),
        bysecond: optionPreparer(origOptions.bysecond),
      };
    },

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
