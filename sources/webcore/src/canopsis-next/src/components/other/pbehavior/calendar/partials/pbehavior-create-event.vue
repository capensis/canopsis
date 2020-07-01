<template lang="pug">
  v-form(@submit.prevent="submit")
    v-layout(row, wrap)
      v-text-field(
        v-model="form.title",
        :label="$t('common.title')",
        name="title"
      )
    v-layout(row, wrap)
      v-text-field(
        v-model="form.rrule",
        :label="$t('common.rrule')",
        name="title"
      )
    v-layout(row, wrap)
      v-flex(xs6)
        v-text-field(
          v-model="form.start_at",
          :label="$t('common.startDate')",
          name="startDate",
          disabled
        )
      v-flex(xs6)
        v-text-field(
          v-model="form.end_at",
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
  props: {
    placeholder: {
      type: Object,
      required: true,
    },
  },
  data() {
    const { placeholder } = this;

    return {
      form: {
        title: placeholder.title,
        rrule: placeholder.pbehavior ? placeholder.pbehavior.rrule : '',
        start_at: placeholder.start.toDate(),
        end_at: placeholder.end.toDate(),
      },
    };
  },
  methods: {
    submit() {
      this.$emit('submit', this.placeholder);
    },
  },
};
</script>
