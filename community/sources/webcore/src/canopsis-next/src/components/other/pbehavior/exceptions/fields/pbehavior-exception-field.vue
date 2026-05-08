<template>
  <v-layout
    class="pbehavior-exception-field"
    column
  >
    <v-layout justify-space-between>
      <v-flex class="pbehavior-exception-field__interval">
        <date-time-splitted-range-picker-field
          v-if="editing"
          :start="value.begin"
          :end="value.end"
          :start-label="$t('common.begin')"
          :end-label="$t('common.end')"
          :start-rules="beginRules"
          :end-rules="endRules"
          :name="datesName"
          :full-day="fullDay"
          :disabled="disabled"
          @update:start="updateField('begin', $event)"
          @update:end="updateField('end', $event)"
        />
        <date-time-splitted-range-picker-text
          v-else
          :start="value.begin"
          :end="value.end"
          :start-label="$t('common.begin')"
          :end-label="$t('common.end')"
          :full-day="fullDay"
        />
      </v-flex>
      <v-flex
        v-if="withType"
        class="pl-2"
      >
        <c-pbehavior-type-field
          v-if="editing"
          v-field="value.type"
          :required="!disabled"
          :name="typeName"
          :disabled="disabled"
          return-object
        />
        <c-pbehavior-type-text
          v-else
          :value="value.type"
        />
      </v-flex>
      <v-flex
        v-if="!disabled"
        class="pbehavior-exception-field__actions"
        shrink
      >
        <v-btn
          :input-value="editing"
          class="btn--editing"
          icon
          fab
          @click="toggleEditing"
        >
          <v-icon
            v-if="editing"
            color="primary"
          >
            check
          </v-icon>
          <v-icon v-else>
            edit
          </v-icon>
        </v-btn>
        <v-btn
          class="v-btn-legacy-m--y"
          color="error"
          icon
          @click="$emit('delete')"
        >
          <v-icon>delete</v-icon>
        </v-btn>
      </v-flex>
    </v-layout>
    <v-layout>
      <v-checkbox
        v-model="fullDay"
        :label="$t('modals.createPbehavior.steps.general.fields.fullDay')"
        :disabled="disabled || !editing"
        class="mt-0"
        color="primary"
        hide-details
      />
    </v-layout>
  </v-layout>
</template>

<script>
import { ref, computed, watch } from 'vue';

import { DATETIME_FORMATS } from '@/constants';

import {
  convertDateToEndOfDayDateObject,
  convertDateToStartOfDayDateObject,
  convertDateToString,
  isEndOfDay,
  isStartOfDay,
} from '@/helpers/date/date';

import { useModelField } from '@/hooks/form/model-field';
import { useValidationChildren } from '@/hooks/validator/validation-children';

import DateTimeSplittedRangePickerField from '@/components/forms/fields/date-time-splitted-range-picker-field.vue';
import DateTimeSplittedRangePickerText from '@/components/forms/fields/date-time-picker/date-time-splitted-range-picker-text.vue';

export default {
  inject: ['$validator'],
  components: {
    DateTimeSplittedRangePickerField,
    DateTimeSplittedRangePickerText,
  },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    withType: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { updateField, updateModel } = useModelField(props, emit);
    const { hasChildrenError, validateChildren } = useValidationChildren();

    const editing = ref(!props.value.type);
    const fullDay = ref(
      isStartOfDay(props.value.begin) && (isEndOfDay(props.value.end) || isStartOfDay(props.value.end)),
    );

    const beginRules = computed(() => ({
      required: true,
      date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
    }));

    const endRules = computed(() => ({
      required: true,
      after: [convertDateToString(props.value.begin, DATETIME_FORMATS.dateTimePicker)],
      date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
    }));

    const nameSuffix = computed(() => (props.value.key ? `-${props.value.key}` : ''));

    const datesName = computed(() => `dates${nameSuffix.value}`);

    const typeName = computed(() => `type${nameSuffix.value}`);

    watch(fullDay, () => {
      updateModel({
        ...props.value,

        begin: convertDateToStartOfDayDateObject(props.value.begin),
        end: convertDateToEndOfDayDateObject(props.value.end),
      });
    });

    const toggleEditing = async () => {
      if (editing.value) {
        await validateChildren();
      }

      if (!hasChildrenError.value) {
        editing.value = !editing.value;
      }
    };

    return {
      editing,
      fullDay,
      beginRules,
      endRules,
      datesName,
      typeName,
      updateField,
      toggleEditing,
    };
  },
};
</script>

<style lang="scss" scoped>
.pbehavior-exception-field {
  &__interval {
    flex-shrink: 0;
  }

  &__actions {
    min-width: 90px;

    .btn--editing {
      height: 36px;
      width: 36px;
    }
  }
}
</style>
