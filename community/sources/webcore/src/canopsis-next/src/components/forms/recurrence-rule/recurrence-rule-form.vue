<template>
  <div class="recurrence-rule-form">
    <c-form-block>
      <c-form-block-row :label="t('recurrenceRule.freq')">
        <recurrence-rule-frequency-field
          :model-value="form.rrule.freq"
          @update:model-value="updateRruleField('freq', $event)"
        />
      </c-form-block-row>

      <template v-if="isFrequencyEnabled">
        <c-form-block-row :label="t('recurrenceRule.count')">
          <c-layout column>
            <recurrence-rule-interval-field
              :model-value="form.rrule"
              @update:model-value="updateFormField('rrule', $event)"
            />

            <c-enabled-field
              v-model="enabledAdvancedFields"
              :label="t('recurrenceRule.tabs.advanced')"
              hide-details
            />

            <v-expand-transition>
              <c-layout v-if="enabledAdvancedFields" column>
                <v-row>
                  <v-col cols="6">
                    <recurrence-rule-weekday-field
                      :model-value="form.rrule.wkst"
                      @update:model-value="updateRruleField('wkst', $event)"
                    />
                  </v-col>

                  <v-col v-if="!isYearlyFrequency" cols="12">
                    <recurrence-rule-weekday-field
                      :model-value="form.rrule.byweekday"
                      chips
                      @update:model-value="updateRruleField('byweekday', $event)"
                    />
                  </v-col>

                  <v-col cols="12">
                    <recurrence-rule-month-field
                      :model-value="form.rrule.bymonth"
                      @update:model-value="updateRruleField('bymonth', $event)"
                    />
                  </v-col>

                  <v-col
                    v-for="field in advancedFields"
                    :key="field.name"
                    cols="6"
                  >
                    <recurrence-rule-advanced-field
                      :model-value="form.rrule[field.name]"
                      :label="t(`recurrenceRule.${field.name}`)"
                      :help-text="t(`recurrenceRule.tooltips.${field.name}`)"
                      :name="field.name"
                      :negative="field.negative"
                      :min="field.min"
                      :max="field.max"
                      @update:model-value="updateRruleField(field.name, $event)"
                    />
                  </v-col>
                </v-row>
              </c-layout>
            </v-expand-transition>
          </c-layout>
        </c-form-block-row>

        <c-form-block-row :label="t('pbehavior.exdates.title')">
          <div class="py-4">
            <pbehavior-exceptions-field
              :exdates="form.exdates"
              :add-button-label="t('pbehavior.exceptions.create')"
              :with-exdate-type="withExdateType"
              @update:exdates="updateFormField('exdates', $event)"
            />
          </div>
        </c-form-block-row>

        <c-form-block-row :label="t('pbehavior.exceptions.title')">
          <div class="py-4">
            <pbehavior-exceptions-list
              v-if="form.exceptions.length"
              :exceptions="form.exceptions"
              @update:exceptions="updateFormField('exceptions', $event)"
            />

            <pbehavior-recurrence-rule-exceptions-list-menu
              :value="form.exceptions"
              @update:value="updateFormField('exceptions', $event)"
            />
          </div>
        </c-form-block-row>

        <c-form-block-row :label="t('common.end')">
          <recurrence-rule-end-field
            :model-value="form.rrule"
            @update:model-value="updateFormField('rrule', $event)"
          />
        </c-form-block-row>

        <c-form-block-row :label="t('common.recurrence')">
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

    <c-alert v-if="recurrenceRuleError" type="error">
      {{ recurrenceRuleError }}
    </c-alert>
  </div>
</template>

<script setup>
import { isNil, isNumber } from 'lodash';
import { RRule } from 'rrule';
import { computed, watch } from 'vue';

import {
  formOptionsToRecurrenceRuleOptions,
  getRecurrenceAdvancedFieldNames,
  recurrenceRuleFormOptionsToRruleBodyString,
} from '@/helpers/entities/shared/recurrence-rule/form';

import { useFormContext, useFormErrorsMap } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import RecurrenceRuleInformation from '@/components/common/reccurence-rule/recurrence-rule-information.vue';
import PbehaviorExceptionsField from '@/components/other/pbehavior/exceptions/fields/pbehavior-exceptions-field.vue';
import PbehaviorRecurrenceRuleExceptionsListMenu from '@/components/other/pbehavior/exceptions/fields/pbehavior-recurrence-rule-exceptions-list-menu.vue';
import PbehaviorExceptionsList from '@/components/other/pbehavior/pbehaviors/partials/pbehavior-exceptions-list.vue';

import RecurrenceRuleAdvancedField from './fields/recurrence-rule-advanced-field.vue';
import RecurrenceRuleEndField from './fields/recurrence-rule-end-field.vue';
import RecurrenceRuleFrequencyField from './fields/recurrence-rule-frequency-field.vue';
import RecurrenceRuleIntervalField from './fields/recurrence-rule-interval-field.vue';
import RecurrenceRuleMonthField from './fields/recurrence-rule-month-field.vue';
import RecurrenceRuleWeekdayField from './fields/recurrence-rule-weekday-field.vue';

const props = defineProps({
  form: {
    type: Object,
    required: true,
  },
  start: {
    type: Date,
    default: undefined,
  },
  withExdateType: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update:form']);

const { t } = useI18n();
const formContext = useFormContext();
const { collectErrors } = useFormErrorsMap();

const isFrequencyEnabled = computed(() => isNumber(props.form.rrule.freq));
const isHourlyFrequency = computed(() => props.form.rrule.freq === RRule.HOURLY);
const isMonthlyFrequency = computed(() => props.form.rrule.freq === RRule.MONTHLY);
const isYearlyFrequency = computed(() => props.form.rrule.freq === RRule.YEARLY);
const recurrenceRuleError = computed(() => collectErrors('recurrenceRule')[0]);

const advancedFields = computed(() => {
  const fields = [
    {
      name: 'bysetpos',
      negative: true,
      min: 1,
      max: 366,
    },
  ];

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

const rruleBodyForInformation = computed(() => recurrenceRuleFormOptionsToRruleBodyString(props.form.rrule),
);

const updateFormField = (field, value) => emit('update:form', {
    ...props.form,
    [field]: value,
  });

const updateRruleField = (field, value) => updateFormField('rrule', {
    ...props.form.rrule,
    [field]: value,
  });

const enabledAdvancedFields = computed({
  get: () => props.form.rrule.advancedEnabled ?? false,
  set: value => updateRruleField('advancedEnabled', value),
});

const validateRecurrenceRrule = () => {
  if (isNil(props.form.rrule.freq)) {
    formContext.setFieldError('recurrenceRule', undefined);

    return;
  }

  try {
    const names = getRecurrenceAdvancedFieldNames(props.form.rrule);
    const rrule = new RRule(formOptionsToRecurrenceRuleOptions(props.form.rrule, names));

    if (!rrule.isFullyConvertibleToText()) {
      formContext.addError('recurrenceRule', t('recurrenceRule.errors.main'));

      return;
    }

    formContext.setFieldError('recurrenceRule', undefined);
  } catch {
    formContext.addError('recurrenceRule', t('recurrenceRule.errors.main'));
  }
};

watch(() => props.form.rrule, validateRecurrenceRrule, { deep: true, immediate: true });
</script>
