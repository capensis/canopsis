<template>
  <div class="recurrence-rule-form">
    <v-layout column>
      <v-layout>
        <v-flex
          class="pr-2"
          xs6
        >
          <recurrence-rule-frequency-field
            :value="form.freq"
            @input="updateFrequency"
          />
          <recurrence-rule-interval-field
            v-if="isFrequencyEnabled"
            v-model="form"
          />
          <recurrence-rule-advanced-repeat-field
            v-if="isMonthlyFrequency"
            v-model="form"
            :start="start"
          />
        </v-flex>
        <v-flex
          class="pl-2"
          xs6
        >
          <recurrence-rule-end-field
            v-if="isFrequencyEnabled"
            v-model="form"
          />
        </v-flex>
      </v-layout>
      <c-collapse-panel
        v-if="isFrequencyEnabled"
        :color="advancedCollapseColor"
        class="my-2"
      >
        <template #header="">
          <span>{{ $t('recurrenceRule.tabs.advanced') }}</span>
        </template>
        <template #actions="">
          <v-icon>$vuetify.icons.expand</v-icon>
        </template>
        <v-layout wrap>
          <v-flex xs6>
            <recurrence-rule-weekday-field v-model="form.wkst" />
          </v-flex>
          <v-flex
            v-if="!isYearlyFrequency"
            xs12
          >
            <recurrence-rule-weekday-field
              v-model="form.byweekday"
              chips
            />
          </v-flex>
          <v-flex xs12>
            <recurrence-rule-month-field v-model="form.bymonth" />
          </v-flex>
          <v-flex
            v-for="(field, index) in advancedFields"
            :key="field.name"
            :class="`${index % 2 ? 'pl' : 'pr'}-2`"
            xs6
          >
            <recurrence-rule-advanced-field
              v-model="form[field.name]"
              :label="$t(`recurrenceRule.${field.name}`)"
              :help-text="$t(`recurrenceRule.tooltips.${field.name}`)"
              :name="field.name"
              :negative="field.negative"
              :min="field.min"
              :max="field.max"
            />
          </v-flex>
        </v-layout>
      </c-collapse-panel>
    </v-layout>
    <recurrence-rule-information
      v-if="isFrequencyEnabled"
      :rrule="recurrenceRuleString"
    />
    <c-alert
      :value="errors.has('recurrenceRule')"
      type="error"
    >
      {{ errors.first('recurrenceRule') }}
    </c-alert>
  </div>
</template>

<script>
import { RRule, rrulestr } from 'rrule';
import { isNull, map } from 'lodash';

import {
  formOptionsToRecurrenceRuleOptions,
  recurrenceRuleToFormOptions,
} from '@/helpers/entities/shared/recurrence-rule/form';

import { formBaseMixin } from '@/mixins/form';

import RecurrenceRuleInformation from '@/components/common/reccurence-rule/recurrence-rule-information.vue';
import RecurrenceRuleAdvancedField from '@/components/forms/recurrence-rule/fields/recurrence-rule-advanced-field.vue';
import RecurrenceRuleEndField from '@/components/forms/recurrence-rule/fields/recurrence-rule-end-field.vue';
import RecurrenceRuleIntervalField from '@/components/forms/recurrence-rule/fields/recurrence-rule-interval-field.vue';

import RecurrenceRuleMonthField from './fields/recurrence-rule-month-field.vue';
import RecurrenceRuleWeekdayField from './fields/recurrence-rule-weekday-field.vue';
import RecurrenceRuleFrequencyField from './fields/recurrence-rule-frequency-field.vue';
import RecurrenceRuleAdvancedRepeatField from './fields/recurrence-rule-advanced-repeat-field.vue';

export default {
  inject: ['$validator', '$system'],
  components: {
    RecurrenceRuleAdvancedRepeatField,
    RecurrenceRuleIntervalField,
    RecurrenceRuleEndField,
    RecurrenceRuleAdvancedField,
    RecurrenceRuleMonthField,
    RecurrenceRuleWeekdayField,
    RecurrenceRuleFrequencyField,
    RecurrenceRuleInformation,
  },
  mixins: [formBaseMixin],
  model: {
    prop: 'rrule',
    event: 'input',
  },
  props: {
    rrule: {
      type: String,
      default: '',
    },
    start: {
      type: Date,
      required: false,
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
      recurrenceRuleObject: recurrenceRule,
      form: recurrenceRuleToFormOptions(recurrenceRule.origOptions),
    };
  },
  computed: {
    isFrequencyEnabled() {
      return !isNull(this.form.freq);
    },

    isHourlyFrequency() {
      return this.form.freq === RRule.HOURLY;
    },

    isWeeklyFrequency() {
      return this.form.freq === RRule.WEEKLY;
    },

    isMonthlyFrequency() {
      return this.form.freq === RRule.MONTHLY;
    },

    isYearlyFrequency() {
      return this.form.freq === RRule.YEARLY;
    },

    advancedFields() {
      const fields = [{
        name: 'bysetpos',
        negative: true,
        min: 1,
        max: 366,
      }];

      if (!this.isMonthlyFrequency) {
        fields.push({
          name: 'byyearday',
          negative: true,
          min: 1,
          max: 366,
        });
      }

      if (!this.isYearlyFrequency) {
        fields.push({
          name: 'bymonthday',
          negative: true,
          min: 1,
          max: 31,
        });
      }

      if (!this.isMonthlyFrequency && !this.isYearlyFrequency) {
        fields.push({
          name: 'byweekno',
          negative: true,
          min: 1,
          max: 53,
        });
      }

      if (this.isHourlyFrequency) {
        fields.push({
          name: 'byhour',
          min: 0,
          max: 23,
        });
      }

      return fields;
    },

    recurrenceRuleString() {
      return this.recurrenceRuleObject.toString();
    },

    advancedCollapseColor() {
      return this.$system.dark ? '#555' : '#e0e0e0';
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
    resetForm() {
      const recurrenceRule = new RRule();

      this.recurrenceRuleObject = recurrenceRule;
      this.form = recurrenceRuleToFormOptions(recurrenceRule.origOptions);
    },

    updateFrequency(frequency) {
      if (!frequency) {
        this.resetForm();
        return;
      }

      this.form.freq = frequency;

      if (!this.isWeeklyFrequency && this.form.byweekday) {
        this.form.byweekday = [];
      }
    },

    /**
     * For each changes in the form we call this function.
     * If RRule isn't valid then add error message to visible RRule field
     * Else remove errors and $emit changes
     */
    changeRecurrenceRuleOption() {
      try {
        this.recurrenceRuleObject = new RRule(
          formOptionsToRecurrenceRuleOptions(
            this.form,
            map(this.advancedFields, 'name'),
          ),
        );

        if (!this.errors.has('recurrenceRule') && !this.recurrenceRuleObject.isFullyConvertibleToText()) {
          this.errors.add({
            field: 'recurrenceRule',
            msg: this.$t('recurrenceRule.errors.main'),
          });
        } else {
          this.errors.remove('recurrenceRule');

          this.updateModel(this.recurrenceRuleString.replace(/.*RRULE:/, ''));
        }
      } catch (err) {
        this.updateModel('');
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
