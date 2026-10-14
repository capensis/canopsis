<template>
  <c-form-block>
    <c-form-block-row :label="t('common.priority')">
      <c-priority-field v-model="priorityField" required>
        <template #append>
          <c-help-icon
            :text="t('broadcastMessage.priorityHelp')"
            max-width="250"
            icon="help"
            left
          />
        </template>
      </c-priority-field>
    </c-form-block-row>

    <c-form-block-row :label="t('common.message')">
      <text-editor-field
        v-model="messageField"
        :label="t('common.message')"
        name="message"
        autofocus
        public
        required
      />
    </c-form-block-row>

    <c-form-block-row :label="t('common.period')">
      <v-row>
        <v-col cols="6">
          <date-time-picker-field
            v-model:value="startField"
            :label="t('common.start')"
            :rules="startRules"
            control-variant="modal"
            name="start"
          >
            <template #append>
              <c-icon>calendar_today</c-icon>
            </template>
          </date-time-picker-field>
        </v-col>
        <v-col cols="6">
          <date-time-picker-field
            v-model:value="endField"
            :label="t('common.end')"
            :rules="endRules"
            control-variant="modal"
            name="end"
          >
            <template #append>
              <c-icon>calendar_today</c-icon>
            </template>
          </date-time-picker-field>
        </v-col>
      </v-row>
    </c-form-block-row>

    <c-form-block-row :label="t('broadcastMessage.closable')" align-center>
      <v-switch v-model="closableField" color="primary" hide-details>
        <template #label>
          {{ t('broadcastMessage.closable') }}
        </template>
        <template #append>
          <c-help-icon
            :text="t('broadcastMessage.closableHelp')"
            max-width="250"
            icon="help"
            left
          />
        </template>
      </v-switch>
    </c-form-block-row>

    <c-form-block-row :label="t('common.color')" align-center>
      <c-color-picker-field v-model="colorField" hide-details />
    </c-form-block-row>
  </c-form-block>
</template>

<script>
import { computed } from 'vue';

import { DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  components: { DateTimePickerField, TextEditorField },

  setup() {
    const { t } = useI18n();

    const {
      priorityField,
      messageField,
      startField,
      endField,
      closableField,
      colorField,
    } = useFormModels({
      priorityField: 'priority',
      messageField: 'message',
      startField: 'start',
      endField: 'end',
      closableField: 'closable',
      colorField: 'color',
    });

    const startRules = computed(() => ({
      required: true,
      date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
    }));

    const endRules = computed(() => ({
      required: true,
      after: [convertDateToString(startField.value, DATETIME_FORMATS.dateTimePicker)],
      date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
    }));

    return {
      t,

      priorityField,
      messageField,
      startField,
      endField,
      closableField,
      colorField,

      startRules,
      endRules,
    };
  },
};
</script>
