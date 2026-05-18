<template>
  <v-layout :class="{ 'gap-3': !blockChild }" column>
    <c-enabled-field
      v-if="withEnabled"
      v-field="form.enabled"
      hide-details
      with-background
    />

    <component :is="blockChild ? 'c-form-block-row' : 'div'" :label="nameLabel" :depth="depth">
      <c-name-field
        v-field="form.name"
        :label="nameLabel"
        :tooltip="nameTooltip"
        autofocus
        required
      />
    </component>

    <component :is="blockChild ? 'div' : 'c-form-block'">
      <c-form-block-row
        v-if="withInherited"
        :label="$t('modals.createPbehavior.steps.general.fields.inherited')"
        :depth="depth"
      >
        <c-enabled-field
          v-field="form.inherited"
          :label="$t('modals.createPbehavior.steps.general.fields.inherited')"
          hide-details
        />
      </c-form-block-row>

      <c-form-block-row
        :label="$t('common.duration')"
        :depth="depth"
      >
        <c-enabled-field
          v-if="withStartOnTrigger"
          :value="form.start_on_trigger"
          :label="$t('modals.createPbehavior.steps.general.fields.startOnTrigger')"
          class="mt-0 mb-1"
          hide-details
          @input="updateStartOnTrigger"
        />

        <c-duration-field
          v-if="form.start_on_trigger"
          v-field="form.duration"
          required
        />
        <template v-else>
          <v-layout class="gap-2 mb-3" align-center>
            <v-flex xs6>
              <c-enabled-field
                v-model="fullDay"
                :label="$t('modals.createPbehavior.steps.general.fields.fullDay')"
                hide-details
              />
            </v-flex>
            <v-flex xs6>
              <c-enabled-field
                v-if="hasPauseType"
                v-model="noEnding"
                :label="$t('modals.createPbehavior.steps.general.fields.noEnding')"
                hide-details
              />
            </v-flex>
          </v-layout>
          <v-layout class="gap-2" align-center>
            <v-flex>
              <date-time-splitted-range-picker-field
                :start="form.tstart"
                :end="form.tstop"
                :start-label="$t('modals.createPbehavior.steps.general.fields.start')"
                :end-label="$t('modals.createPbehavior.steps.general.fields.stop')"
                :start-rules="tstartRules"
                :end-rules="tstopRules"
                :end-min="tstopMin"
                :end-max="tstopMax"
                :no-ending="noEnding"
                :full-day="fullDay"
                @update:start="updateField('tstart', $event)"
                @update:end="updateField('tstop', $event)"
              />
            </v-flex>
            <v-flex v-if="!noTimezone">
              <c-timezone-field
                v-field="form.timezone"
                server
              />
            </v-flex>
          </v-layout>
        </template>
      </c-form-block-row>

      <c-form-block-row :label="'Pbehavior reason and type'" :depth="depth">
        <v-layout>
          <v-flex xs6>
            <c-pbehavior-reason-field
              v-field="form.reason"
              class="mr-2"
              required
              return-object
            />
          </v-flex>
          <v-flex xs6>
            <c-pbehavior-type-field
              v-field="form.type"
              class="ml-2"
              required
              return-object
            />
          </v-flex>
        </v-layout>
      </c-form-block-row>

      <c-form-block-row :label="$t('common.color')" :depth="depth">
        <c-enabled-color-picker-field
          v-field="form.color"
          :label="$t('modals.createPbehavior.steps.color.label')"
          row
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('modals.createPbehavior.steps.rrule.title')" :depth="depth" indented>
        <pbehavior-recurrence-rule-field v-field="form" with-exdate-type />
      </c-form-block-row>

      <c-form-block-row
        v-if="!noComments"
        :label="$tc('common.comment', 2)"
        :depth="depth"
        indented
      >
        <pbehavior-comments-field v-field="form.comments" />
      </c-form-block-row>
    </component>
  </v-layout>
</template>

<script>
import { computed, onMounted, ref, watch } from 'vue';

import { MAX_PBEHAVIOR_DATES_DIFF_YEARS } from '@/config';
import { DATETIME_FORMATS, PBEHAVIOR_TYPE_TYPES, TIME_UNITS } from '@/constants';

import {
  isStartOfDay,
  isEndOfDay,
  addUnitToDate,
  getNowTimestamp,
  convertDateToString,
  convertDateToTimestamp,
  convertDateToDateObject,
  convertDateToStartOfDayDateObject,
  convertDateToEndOfUnitDateObject,
  convertDateToEndOfDayDateObject,
} from '@/helpers/date/date';

import { useModelField } from '@/hooks/form/model-field';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';

import DateTimeSplittedRangePickerField from '@/components/forms/fields/date-time-splitted-range-picker-field.vue';
import PbehaviorRecurrenceRule from '@/components/modals/pbehavior/pbehavior-recurrence-rule.vue';

import PbehaviorCommentsField from '../fields/pbehavior-comments-field.vue';
import PbehaviorRecurrenceRuleField from '../fields/pbehavior-recurrence-rule-field.vue';

export default {
  inject: ['$validator'],
  components: {
    PbehaviorCommentsField,
    DateTimeSplittedRangePickerField,
    PbehaviorRecurrenceRule,
    PbehaviorRecurrenceRuleField,
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
    noEnabled: {
      type: Boolean,
      default: false,
    },
    noComments: {
      type: Boolean,
      default: false,
    },
    noTimezone: {
      type: Boolean,
      default: false,
    },
    withStartOnTrigger: {
      type: Boolean,
      default: false,
    },
    withInherited: {
      type: Boolean,
      default: false,
    },
    withEnabled: {
      type: Boolean,
      default: false,
    },
    nameLabel: {
      type: String,
      required: false,
    },
    nameTooltip: {
      type: String,
      required: false,
    },
    depth: {
      type: Number,
      default: 0,
    },
    blockChild: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { updateModel, updateField } = useModelField(props, emit);
    const { fetchPbehaviorTypesFieldList } = usePbehaviorType();

    const initialNoEnding = Boolean(props.form.tstart && !props.form.tstop);
    const noEnding = ref(initialNoEnding);
    const fullDay = ref(
      isStartOfDay(props.form.tstart) && (initialNoEnding || isEndOfDay(props.form.tstop)),
    );

    const hasPauseType = computed(() => props.form.type?.type === PBEHAVIOR_TYPE_TYPES.pause);

    const tstartRules = computed(() => ({
      required: true,
      date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
    }));

    const tstopRules = computed(() => {
      const rules = { required: !hasPauseType.value };

      if (props.form.tstart) {
        rules.after = [convertDateToString(props.form.tstart, DATETIME_FORMATS.dateTimePicker)];
        rules.date_format = DATETIME_FORMATS.veeValidateDateTimeFormat;
      }

      return rules;
    });

    const tstopMin = computed(() => {
      const nowTimestamp = getNowTimestamp();
      const startTimestamp = convertDateToTimestamp(props.form.tstart);

      return convertDateToString(
        Math.min(nowTimestamp, startTimestamp),
        DATETIME_FORMATS.vuetifyDatePicker,
      );
    });

    const tstopMax = computed(() => convertDateToString(
      addUnitToDate(props.form.tstart, MAX_PBEHAVIOR_DATES_DIFF_YEARS, TIME_UNITS.year),
      DATETIME_FORMATS.vuetifyDatePicker,
    ));

    watch(noEnding, (isNoEnding) => {
      const { tstart } = props.form;

      if (isNoEnding) {
        updateField('tstop', null);
      } else if (tstart) {
        const unit = fullDay.value ? 'day' : 'hour';

        const tstop = addUnitToDate(tstart, 1, unit);
        const tstopDate = fullDay.value
          ? convertDateToEndOfUnitDateObject(tstop, unit)
          : convertDateToDateObject(tstop);

        updateField('tstop', tstopDate);
      }
    });

    watch(fullDay, () => {
      const { tstart, tstop } = props.form;

      if (tstart) {
        updateModel({
          ...props.form,

          tstart: convertDateToStartOfDayDateObject(tstart),
          tstop: !noEnding.value && tstop ? convertDateToEndOfDayDateObject(tstop) : tstop,
        });
      }
    });

    watch(hasPauseType, (value) => {
      if (!value) {
        noEnding.value = false;
      }
    });

    /**
     * Syncs form state when "start on trigger" is toggled. Enabling clears scheduled window flags
     * and removes start/stop; disabling drops duration and turns start_on_trigger off.
     *
     * @param {boolean} value - Whether the pbehavior should start only when the trigger fires.
     */
    const updateStartOnTrigger = (value) => {
      if (value) {
        fullDay.value = false;
        noEnding.value = false;

        updateModel({
          ...props.form,

          start_on_trigger: true,
          tstart: null,
          tstop: null,
        });
      } else {
        const { duration, ...form } = props.form;

        updateModel({
          ...form,

          start_on_trigger: false,
        });
      }
    };

    onMounted(fetchPbehaviorTypesFieldList);

    return {
      noEnding,
      fullDay,
      hasPauseType,
      tstartRules,
      tstopRules,
      tstopMin,
      tstopMax,
      updateField,
      updateStartOnTrigger,
    };
  },
};
</script>
