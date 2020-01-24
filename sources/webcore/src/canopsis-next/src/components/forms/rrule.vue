<template lang="pug">
  div
    v-layout(row)
      v-switch(
        data-test="pbehaviorRuleSwitcher",
        v-model="showRRule",
        :label="$t('modals.createPbehavior.steps.rrule.fields.rRuleQuestion')",
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
                  :label="$t('rRule.fields.freq')",
                  v-model="form.rRuleOptions.freq",
                  :items="frequencies",
                  @input="changeRRuleOption"
                )
              div(data-test="pbehaviorByWeekDay")
                v-select(
                  :label="$t('rRule.fields.byweekday')",
                  v-model="form.rRuleOptions.byweekday",
                  :items="weekDays",
                  multiple,
                  chips,
                  @input="changeRRuleOption"
                )
              div
                v-text-field(
                  data-test="pbehaviorRepeat",
                  type="number",
                  :label="$t('rRule.fields.count')",
                  v-model="form.rRuleOptions.count",
                  @input="changeRRuleOption"
                )
              div
                v-text-field(
                  data-test="pbehaviorInterval",
                  type="number",
                  :label="$t('rRule.fields.interval')",
                  v-model="form.rRuleOptions.interval",
                  @input="changeRRuleOption"
                )
          v-tab-item(value="advanced")
            div
              div(data-test="pbehaviorWeekStart")
                v-select(
                  :label="$t('rRule.fields.wkst')",
                  v-model="form.rRuleOptions.wkst",
                  :items="weekDays",
                  @input="changeRRuleOption"
                )
              div(data-test="pbehaviorByMonth")
                v-select(
                  :label="$t('rRule.fields.bymonth')",
                  v-model="form.rRuleOptions.bymonth",
                  :items="months",
                  @input="changeRRuleOption",
                  multiple,
                  chips
                )
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      data-test="pbehaviorBySetPos",
                      :label="$t('rRule.fields.bysetpos.label')",
                      v-model="form.advancedRRuleOptions.bysetpos",
                      :error-messages="errors.collect('bysetpos')",
                      :hint="$t('rRule.advancedHint')",
                      @input="changeRRuleOption",
                      v-validate="{ regex: advancedFieldRegex }",
                      data-vv-name="bysetpos",
                      persistent-hint
                    )
                  span {{ $t('rRule.fields.bysetpos.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      data-test="pbehaviorByMonthDay",
                      :label="$t('rRule.fields.bymonthday.label')",
                      v-model="form.advancedRRuleOptions.bymonthday",
                      :error-messages="errors.collect('bymonthday')",
                      :hint="$t('rRule.advancedHint')",
                      @input="changeRRuleOption",
                      v-validate="{ regex: advancedFieldRegex }",
                      data-vv-name="bymonthday",
                      persistent-hint
                    )
                  span {{ $t('rRule.fields.bymonthday.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      data-test="pbehaviorByYearDay",
                      :label="$t('rRule.fields.byyearday.label')",
                      v-model="form.advancedRRuleOptions.byyearday",
                      :error-messages="errors.collect('byyearday')",
                      :hint="$t('rRule.advancedHint')",
                      @input="changeRRuleOption",
                      v-validate="{ regex: advancedFieldRegex }",
                      data-vv-name="byyearday",
                      persistent-hint
                    )
                  span {{ $t('rRule.fields.byyearday.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      data-test="pbehaviorByWeekNo",
                      :label="$t('rRule.fields.byweekno.label')",
                      v-model="form.advancedRRuleOptions.byweekno",
                      :error-messages="errors.collect('byweekno')",
                      :hint="$t('rRule.advancedHint')",
                      @input="changeRRuleOption",
                      v-validate="{ regex: advancedFieldRegex }",
                      data-vv-name="byweekno",
                      persistent-hint
                    )
                  span {{ $t('rRule.fields.byweekno.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      data-test="pbehaviorByHour",
                      :label="$t('rRule.fields.byhour.label')",
                      v-model="form.advancedRRuleOptions.byhour",
                      :error-messages="errors.collect('byhour')",
                      :hint="$t('rRule.advancedHint')",
                      @input="changeRRuleOption",
                      v-validate="{ regex: advancedFieldRegex }",
                      data-vv-name="byhour",
                      persistent-hint
                    )
                  span {{ $t('rRule.fields.byhour.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      data-test="pbehaviorByMinute",
                      :label="$t('rRule.fields.byminute.label')",
                      v-model="form.advancedRRuleOptions.byminute",
                      :error-messages="errors.collect('byminute')",
                      :hint="$t('rRule.advancedHint')",
                      @input="changeRRuleOption",
                      v-validate="{ regex: advancedFieldRegex }",
                      data-vv-name="byminute",
                      persistent-hint
                    )
                  span {{ $t('rRule.fields.byminute.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                      data-test="pbehaviorBySecond",
                      :label="$t('rRule.fields.bysecond.label')",
                      v-model="form.advancedRRuleOptions.bysecond",
                      :error-messages="errors.collect('bysecond')",
                      :hint="$t('rRule.advancedHint')",
                      @input="changeRRuleOption",
                      v-validate="{ regex: advancedFieldRegex }",
                      data-vv-name="bysecond",
                      persistent-hint
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
  filters: {
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

    formOptionsToRRuleOptions(options) {
      return pickBy(options, v => v !== '');
    },

    formAdvancedOptionsToRRuleOptions(options) {
      return mapValues(options, o => o.split(',').filter(v => v));
    },
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
  inject: ['$validator'],
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
        rRuleOptions: this.$options.filters.rRuleToFormOptions(rRule),
        advancedRRuleOptions: this.$options.filters.rRuleToFormAdvancedOptions(rRule),
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
    /**
     * For each changes in the form we call this function.
     * If RRule isn't valid then add error message to visible RRule field
     * Else remove errors and $emit changes
     */
    changeRRuleOption() {
      try {
        this.rRuleObject = new RRule({
          ...this.$options.filters.formOptionsToRRuleOptions(this.form.rRuleOptions),
          ...this.$options.filters.formAdvancedOptionsToRRuleOptions(this.form.advancedRRuleOptions),
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
