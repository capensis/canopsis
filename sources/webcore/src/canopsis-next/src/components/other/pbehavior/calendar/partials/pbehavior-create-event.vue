<template lang="pug">
  div
    v-layout(row, wrap)
      v-text-field(
        v-model="form.title",
        :label="$t('common.title')",
        name="title"
      )
    v-layout(row, wrap)
      v-flex(xs6)
        date-time-picker-text-field(
          v-model="form.startDate",
          :label="$t('common.startDate')",
          name="startDate"
        )
      v-flex(xs6)
        date-time-picker-text-field(
          v-model="form.endDate",
          :label="$t('common.endDate')",
          name="endDate"
        )
    v-layout(row, justify-end)
      v-btn.mr-0.mb-0(
        depressed,
        flat,
        @click="close"
      ) {{ $t('common.cancel') }}
      v-btn.mr-0.mb-0.primary.white--text(type="submit") {{ $t('common.submit') }}
</template>

<script>
import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

export default {
  components: { DateTimePickerTextField },
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    close: {
      type: Function,
      required: true,
    },
  },
};
</script>
