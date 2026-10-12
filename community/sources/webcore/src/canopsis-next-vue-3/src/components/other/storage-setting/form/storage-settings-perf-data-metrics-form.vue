<template>
  <c-information-block :title="$t('storageSetting.perfDataMetrics.title')">
    <c-enabled-duration-field
      :duration="form.delete_after"
      :label="$t('storageSetting.perfDataMetrics.deleteAfter')"
      :suffix="$t('common.after')"
      name="perf_data_metrics.delete_after"
      switcher
      hide-value-on-false
      @update:duration="updateField('delete_after', $event)"
    />
  </c-information-block>
</template>

<script>
import { useModelField } from '@/hooks/form';

export default {
  props: {
    form: {
      type: Object,
      required: true,
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    return { updateField };
  },
};
</script>
