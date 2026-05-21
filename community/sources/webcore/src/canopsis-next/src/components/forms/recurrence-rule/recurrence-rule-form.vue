<template>
  <div class="recurrence-rule-form">
    <c-form-block>
      <c-form-block-row :label="$t('recurrenceRule.freq')">
        <recurrence-rule-frequency-field
          :value="form.rrule.freq"
          @input="updateFrequency"
        />
      </c-form-block-row>
      <template v-if="isFrequencyEnabled">
        <c-form-block-row :label="$t('recurrenceRule.count')">
          <v-layout column>
            <recurrence-rule-interval-field v-field="form.rrule" />

            <c-enabled-field v-model="enabledAdvancedFields" :label="$t('recurrenceRule.tabs.advanced')" />

            <v-expand-transition>
              <v-layout v-if="enabledAdvancedFields" column>
                <v-layout wrap>
                  <v-flex xs6>
                    <recurrence-rule-weekday-field v-field="form.rrule.wkst" />
                  </v-flex>
                  <v-flex
                    v-if="!isYearlyFrequency"
                    xs12
                  >
                    <recurrence-rule-weekday-field
                      v-field="form.rrule.byweekday"
                      chips
                    />
                  </v-flex>
                  <v-flex xs12>
                    <recurrence-rule-month-field v-field="form.rrule.bymonth" />
                  </v-flex>
                  <v-flex
                    v-for="(field, index) in advancedFields"
                    :key="field.name"
                    :class="`${index % 2 ? 'pl' : 'pr'}-2`"
                    xs6
                  >
                    <recurrence-rule-advanced-field
                      v-field="form.rrule[field.name]"
                      :label="$t(`recurrenceRule.${field.name}`)"
                      :help-text="$t(`recurrenceRule.tooltips.${field.name}`)"
                      :name="field.name"
                      :negative="field.negative"
                      :min="field.min"
                      :max="field.max"
                    />
                  </v-flex>
                </v-layout>
              </v-layout>
            </v-expand-transition>
          </v-layout>
        </c-form-block-row>
        <c-form-block-row :label="$t('pbehavior.exdates.title')">
          <div class="py-4">
            <pbehavior-exceptions-field
              v-field="form.exdates"
              :add-button-label="$t('pbehavior.exceptions.create')"
              :with-exdate-type="withExdateType"
            />
          </div>
        </c-form-block-row>
        <c-form-block-row :label="$t('pbehavior.exceptions.title')">
          <div class="py-4">
            <pbehavior-exceptions-list
              v-if="form.exceptions.length"
              v-field="form.exceptions"
            />
            <pbehavior-recurrence-rule-exceptions-list-menu v-field="form.exceptions" />
          </div>
        </c-form-block-row>
        <c-form-block-row :label="$t('common.end')">
          <recurrence-rule-end-field v-field="form.rrule" />
        </c-form-block-row>
        <c-form-block-row v-if="isFrequencyEnabled" :label="$t('common.recurrence')">
          <div class="py-4">
            <recurrence-rule-information
              :rrule="rruleBodyForInformation"
              :exdates="form.exdates"
              :exceptions="form.exceptions"
            />
          </div>
        </c-form-block-row>
      </template>
    </c-form-block>
    <c-alert
      :value="errors.has('recurrenceRule')"
      type="error"
    >
      {{ errors.first('recurrenceRule') }}
    </c-alert>
  </div>
</template>

<script>
import { ref, computed, watch, onMounted } from 'vue';
import { RRule } from 'rrule';
import { isNumber, isEmpty } from 'lodash';

import {
  emptyRecurrenceRuleFormOptions,
  formOptionsToRecurrenceRuleOptions,
  getRecurrenceAdvancedFieldNames,
  recurrenceRuleFormOptionsToRruleBodyString,
} from '@/helpers/entities/shared/recurrence-rule/form';
import { convertDateToStartOfDayDateObject, convertDateToEndOfDayDateObject } from '@/helpers/date/date';
import { uid } from '@/helpers/uid';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';

import RecurrenceRuleInformation from '@/components/common/reccurence-rule/recurrence-rule-information.vue';
import RecurrenceRuleAdvancedField from '@/components/forms/recurrence-rule/fields/recurrence-rule-advanced-field.vue';
import RecurrenceRuleEndField from '@/components/forms/recurrence-rule/fields/recurrence-rule-end-field.vue';
import RecurrenceRuleIntervalField from '@/components/forms/recurrence-rule/fields/recurrence-rule-interval-field.vue';
import PbehaviorExceptionsList from '@/components/other/pbehavior/pbehaviors/partials/pbehavior-exceptions-list.vue';
import PbehaviorExceptionsField from '@/components/other/pbehavior/exceptions/fields/pbehavior-exceptions-field.vue';
import PbehaviorRecurrenceRuleExceptionsListMenu from '@/components/other/pbehavior/exceptions/fields/pbehavior-recurrence-rule-exceptions-list-menu.vue';

import RecurrenceRuleMonthField from './fields/recurrence-rule-month-field.vue';
import RecurrenceRuleWeekdayField from './fields/recurrence-rule-weekday-field.vue';
import RecurrenceRuleFrequencyField from './fields/recurrence-rule-frequency-field.vue';

export default {
  inject: ['$validator'],
  components: {
    RecurrenceRuleIntervalField,
    RecurrenceRuleEndField,
    RecurrenceRuleAdvancedField,
    RecurrenceRuleMonthField,
    RecurrenceRuleWeekdayField,
    RecurrenceRuleFrequencyField,
    RecurrenceRuleInformation,
    PbehaviorExceptionsList,
    PbehaviorExceptionsField,
    PbehaviorRecurrenceRuleExceptionsListMenu,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    start: {
      type: Date,
      required: false,
    },
    withExdateType: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const validator = useValidator();
    const { updateField } = useModelField(props, emit);

    const isFrequencyEnabled = computed(() => isNumber(props.form.rrule.freq));

    const isHourlyFrequency = computed(() => props.form.rrule.freq === RRule.HOURLY);

    const isWeeklyFrequency = computed(() => props.form.rrule.freq === RRule.WEEKLY);

    const isMonthlyFrequency = computed(() => props.form.rrule.freq === RRule.MONTHLY);

    const isYearlyFrequency = computed(() => props.form.rrule.freq === RRule.YEARLY);

    const advancedFields = computed(() => {
      const fields = [{
        name: 'bysetpos',
        negative: true,
        min: 1,
        max: 366,
      }];

      if (!isMonthlyFrequency.value) {
        fields.push({
          name: 'byyearday',
          negative: true,
          min: 1,
          max: 366,
        });
      }

      if (!isYearlyFrequency.value) {
        fields.push({
          name: 'bymonthday',
          negative: true,
          min: 1,
          max: 31,
        });
      }

      if (!isMonthlyFrequency.value && !isYearlyFrequency.value) {
        fields.push({
          name: 'byweekno',
          negative: true,
          min: 1,
          max: 53,
        });
      }

      if (isHourlyFrequency.value) {
        fields.push({
          name: 'byhour',
          min: 0,
          max: 23,
        });
      }

      return fields;
    });

    const hasChangedAdvancedFields = () => {
      const { rrule } = props.form;

      return ![

        ...advancedFields.value.map(field => rrule[field.name]),
        rrule.wkst,
        rrule.byweekday,
        rrule.bymonth,
      ].every(isEmpty);
    };

    const enabledAdvancedFields = ref(hasChangedAdvancedFields());

    const rruleBodyForInformation = computed(() => (
      recurrenceRuleFormOptionsToRruleBodyString(props.form.rrule)
    ));

    const validateRecurrenceRrule = () => {
      if (props.form.rrule.freq == null) {
        validator.errors.remove('recurrenceRule');

        return;
      }

      try {
        const names = getRecurrenceAdvancedFieldNames(props.form.rrule);
        const rruleInst = new RRule(formOptionsToRecurrenceRuleOptions(props.form.rrule, names));

        if (!rruleInst.isFullyConvertibleToText()) {
          validator.errors.add({
            field: 'recurrenceRule',
            msg: t('recurrenceRule.errors.main'),
          });

          return;
        }

        validator.errors.remove('recurrenceRule');
      } catch {
        validator.errors.add({
          field: 'recurrenceRule',
          msg: t('recurrenceRule.errors.main'),
        });
      }
    };

    watch(() => props.form.rrule, () => {
      validateRecurrenceRrule();
      enabledAdvancedFields.value = hasChangedAdvancedFields();
    }, { deep: true });

    onMounted(() => {
      validateRecurrenceRrule();
      enabledAdvancedFields.value = hasChangedAdvancedFields();
    });

    const updateFrequency = (frequency) => {
      if (frequency === null) {
        updateField('rrule', emptyRecurrenceRuleFormOptions());

        return;
      }

      const { rrule } = props.form;
      const nextRrule = { ...rrule, freq: frequency };

      if (frequency !== RRule.WEEKLY && rrule.byweekday?.length) {
        nextRrule.byweekday = [];
      }

      updateField('rrule', nextRrule);
    };

    const addExdate = () => {
      updateField('exdates', [
        ...(props.form.exdates ?? []),
        {
          key: uid(),
          begin: convertDateToStartOfDayDateObject(),
          end: convertDateToEndOfDayDateObject(),
          type: '',
        },
      ]);
    };

    return {
      validator,
      updateField,
      isFrequencyEnabled,
      isHourlyFrequency,
      isWeeklyFrequency,
      isMonthlyFrequency,
      isYearlyFrequency,
      advancedFields,
      enabledAdvancedFields,
      rruleBodyForInformation,
      updateFrequency,
      addExdate,
    };
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
