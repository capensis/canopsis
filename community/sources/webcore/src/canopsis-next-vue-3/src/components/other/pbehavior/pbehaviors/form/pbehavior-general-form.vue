<template>
  <c-layout :gap="blockChild ? undefined : '3'" column>
    <c-enabled-field
      v-if="withEnabled"
      :model-value="form.enabled"
      hide-details
      with-background
      @update:model-value="updateField('enabled', $event)"
    />

    <component :is="blockChild ? 'c-form-block-row' : 'div'" :label="nameLabel" :depth="depth">
      <c-name-field
        :model-value="form.name"
        :label="nameLabel"
        :tooltip="nameTooltip"
        :name="nameFieldName"
        autofocus
        required
        @update:model-value="updateField('name', $event)"
      />
    </component>

    <component :is="blockChild ? 'div' : 'c-form-block'">
      <c-form-block-row
        v-if="withInherited"
        :label="$t('modals.createPbehavior.steps.general.fields.inherited')"
        :depth="depth"
      >
        <c-enabled-field
          :model-value="form.inherited"
          :label="$t('modals.createPbehavior.steps.general.fields.inherited')"
          hide-details
          @update:model-value="updateField('inherited', $event)"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.duration')" :depth="depth">
        <c-enabled-field
          v-if="withStartOnTrigger"
          :model-value="form.start_on_trigger"
          :label="$t('modals.createPbehavior.steps.general.fields.startOnTrigger')"
          class="mt-0 mb-1"
          hide-details
          @update:model-value="updateStartOnTrigger"
        />

        <c-duration-field
          v-if="form.start_on_trigger"
          :duration="form.duration"
          :name="durationFieldName"
          required
          @update:duration="updateField('duration', $event)"
        />

        <template v-else>
          <c-layout class="mb-3" gap="2" align-center>
            <div class="pbehavior-general-form__half">
              <c-enabled-field
                v-model="fullDay"
                :label="$t('modals.createPbehavior.steps.general.fields.fullDay')"
                hide-details
              />
            </div>

            <div class="pbehavior-general-form__half">
              <c-enabled-field
                v-if="hasPauseType"
                v-model="noEnding"
                :label="$t('modals.createPbehavior.steps.general.fields.noEnding')"
                hide-details
              />
            </div>
          </c-layout>

          <c-layout gap="2" align-center>
            <div class="pbehavior-general-form__date-range">
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
                :name="dateTimeSplittedRangePickerFieldName"
                @update:start="updateField('tstart', $event)"
                @update:end="updateField('tstop', $event)"
              />
            </div>

            <div v-if="!noTimezone" class="pbehavior-general-form__timezone">
              <c-timezone-field
                :model-value="form.timezone"
                server
                @update:model-value="updateField('timezone', $event)"
              />
            </div>
          </c-layout>
        </template>
      </c-form-block-row>

      <c-form-block-row :label="'Pbehavior reason and type'" :depth="depth">
        <c-layout gap="3">
          <div class="pbehavior-general-form__half">
            <c-pbehavior-reason-field
              :model-value="form.reason"
              class="mr-2"
              required
              return-object
              @update:model-value="updateField('reason', $event)"
            />
          </div>

          <div class="pbehavior-general-form__half">
            <c-pbehavior-type-field
              :model-value="form.type"
              :name="pbehaviorTypeFieldName"
              class="ml-2"
              required
              return-object
              @update:model-value="updateField('type', $event)"
            />
          </div>
        </c-layout>
      </c-form-block-row>

      <c-form-block-row :label="$t('common.color')" :depth="depth">
        <c-enabled-color-picker-field
          :model-value="form.color"
          :label="$t('modals.createPbehavior.steps.color.label')"
          justify-end
          @update:model-value="updateField('color', $event)"
        />
      </c-form-block-row>

      <c-form-block-row
        :label="$t('modals.createPbehavior.steps.rrule.title')"
        :depth="depth"
        indented
      >
        <pbehavior-recurrence-rule-field :form="form" with-exdate-type @update:form="updateModel" />
      </c-form-block-row>

      <c-form-block-row
        v-if="!noComments"
        :label="$t('common.comment', 2)"
        :depth="depth"
        indented
      >
        <pbehavior-comments-field
          :comments="form.comments"
          @update:comments="updateField('comments', $event)"
        />
      </c-form-block-row>
    </component>
  </c-layout>
</template>

<script>
import { computed, onMounted, ref, watch } from 'vue';

import { MAX_PBEHAVIOR_DATES_DIFF_YEARS } from '@/config';
import { DATETIME_FORMATS, PBEHAVIOR_TYPE_TYPES, TIME_UNITS } from '@/constants';

import {
  addUnitToDate,
  convertDateToDateObject,
  convertDateToEndOfDayDateObject,
  convertDateToEndOfUnitDateObject,
  convertDateToStartOfDayDateObject,
  convertDateToString,
  convertDateToTimestamp,
  getNowTimestamp,
  isEndOfDay,
  isStartOfDay,
} from '@/helpers/date/date';

import { useModelField } from '@/hooks/form/model-field';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';

import DateTimeSplittedRangePickerField from '@/components/forms/fields/date-time-splitted-range-picker-field.vue';

import PbehaviorCommentsField from '../fields/pbehavior-comments-field.vue';
import PbehaviorRecurrenceRuleField from '../fields/pbehavior-recurrence-rule-field.vue';

export default {
  components: {
    PbehaviorCommentsField,
    DateTimeSplittedRangePickerField,
    PbehaviorRecurrenceRuleField,
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
    name: {
      type: String,
      default: '',
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

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateModel, updateField } = useModelField(props, emit);
    const { fetchPbehaviorTypesFieldList } = usePbehaviorType();

    const getFieldName = field => (props.name ? `${props.name}.${field}` : field);

    const nameFieldName = computed(() => getFieldName('name'));
    const durationFieldName = computed(() => getFieldName('duration'));
    const pbehaviorTypeFieldName = computed(() => getFieldName('type'));
    const dateTimeSplittedRangePickerFieldName = computed(() => getFieldName('tstart'));

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
      ),
    );

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
      nameFieldName,
      durationFieldName,
      pbehaviorTypeFieldName,
      dateTimeSplittedRangePickerFieldName,
      tstartRules,
      tstopRules,
      tstopMin,
      tstopMax,
      updateModel,
      updateField,
      updateStartOnTrigger,
    };
  },
};
</script>

<style lang="scss" scoped>
.pbehavior-general-form {
  &__half {
    flex: 1 1 50%;
    min-width: 0;
  }

  &__date-range {
    flex: 1 1 auto;
    min-width: 0;
  }

  &__timezone {
    flex: 0 0 260px;
    min-width: 0;
  }
}
</style>
