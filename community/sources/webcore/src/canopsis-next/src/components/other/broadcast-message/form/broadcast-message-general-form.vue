<template>
  <div>
    <v-layout>
      <v-flex xs12>
        <text-editor-field
          v-field="form.message"
          v-validate="'required'"
          :label="$t('common.message')"
          :error-messages="errors.collect('message')"
          name="message"
          autofocus
          public
        />
      </v-flex>
    </v-layout>
    <v-layout>
      <c-color-picker-field v-field="form.color" />
    </v-layout>
    <v-layout class="gap-4">
      <v-flex xs6>
        <date-time-picker-field
          v-validate="startRules"
          :value="form.start"
          :label="$t('common.start')"
          :error-message="errors.collect('start')"
          name="start"
          @input="updateField('start', $event)"
        >
          <template #append-outer>
            <v-btn icon>
              <v-icon>
                calendar_today
              </v-icon>
            </v-btn>
          </template>
        </date-time-picker-field>
      </v-flex>
      <v-flex xs6>
        <date-time-picker-field
          v-validate="endRules"
          :value="form.end"
          :label="$t('common.end')"
          :error-message="errors.collect('end')"
          name="end"
          @input="updateField('end', $event)"
        >
          <template #append-outer>
            <v-btn icon>
              <v-icon>
                calendar_today
              </v-icon>
            </v-btn>
          </template>
        </date-time-picker-field>
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
import { computed } from 'vue';

import { DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import { useModelField } from '@/hooks/form/model-field';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  inject: ['$validator'],
  components: { TextEditorField, DateTimePickerField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const startRules = computed(() => ({
      required: true,
      date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
    }));

    const endRules = computed(() => ({
      required: true,
      after: [convertDateToString(props.form.start, DATETIME_FORMATS.dateTimePicker)],
      date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
    }));

    return {
      startRules,
      endRules,

      /**
       * We need to use updateField function instead of v-field to avoid problem with Date objects
       */
      updateField,
    };
  },
};
</script>
