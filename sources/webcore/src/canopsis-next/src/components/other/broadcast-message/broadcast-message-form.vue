<template lang="pug">
  div.mt-4
    v-layout(row)
      v-textarea(
        v-field="form.message",
        v-validate="'required'",
        :label="$t('common.message')",
        :error-messages="errors.collect('message')",
        name="message"
      )
    v-layout(row)
      v-btn.ml-0.mt-2(
        :style="{ backgroundColor: form.color }",
        @click="showColorPickerModal",
        dark
      ) {{ $t('modals.createBroadcastMessage.buttons.selectColor') }}
    v-layout(row)
      v-switch(v-field="form.enabled", :label="$t('common.enabled')")
    v-layout(row)
      date-time-picker-field(
        v-validate="startRules",
        :value="form.start",
        :label="$t('common.start')",
        :error-message="errors.collect('start')",
        name="start",
        @input="updateField('start', $event)"
      )
    v-layout(row)
      date-time-picker-field(
        v-validate="endRules",
        :value="form.end",
        :label="$t('common.end')",
        :error-message="errors.collect('end')",
        name="end",
        @input="updateField('end', $event)"
      )
</template>

<script>
import moment from 'moment';

import { DATETIME_FORMATS, MODALS } from '@/constants';

import formMixin from '@/mixins/form';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimePickerField },
  mixins: [formMixin],
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
  computed: {
    startRules() {
      return {
        required: true,
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    endRules() {
      return {
        required: true,
        after: [moment(this.form.start).format(DATETIME_FORMATS.dateTimePicker)],
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },
  },
  methods: {
    showColorPickerModal() {
      this.$modals.show({
        name: MODALS.colorPicker,
        config: {
          color: this.form.color,
          action: color => this.updateField('color', color),
        },
      });
    },
  },
};
</script>
