<template lang="pug">
  v-form(@submit.prevent="$emit('submit', form)")
    v-layout(row, wrap)
      v-text-field(
        v-model="form.title",
        :label="$t('common.title')",
        name="title",
        disabled
      )
    v-layout(row, wrap)
      v-flex(xs6)
        v-text-field(
          v-model="form.start",
          :label="$t('common.startDate')",
          name="startDate",
          disabled
        )
      v-flex(xs6)
        v-text-field(
          v-model="form.end",
          :label="$t('common.endDate')",
          name="endDate",
          disabled
        )
    v-layout(row, justify-end)
      v-btn.mr-0.mb-0(
        depressed,
        flat,
        @click="$emit('close')"
      ) {{ $t('common.cancel') }}
      v-btn.mr-0.mb-0.primary.white--text(type="submit") {{ $t('common.submit') }}
</template>

<script>
export default {
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
  },
};
</script>
