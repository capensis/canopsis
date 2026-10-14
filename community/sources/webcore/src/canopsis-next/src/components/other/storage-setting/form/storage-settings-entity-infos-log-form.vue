<template>
  <c-information-block :title="$t('storageSetting.entityInfosLog.title')">
    <c-enabled-duration-field
      :duration="form.delete_after"
      :label="$t('storageSetting.entityInfosLog.deleteAfter')"
      :name="name"
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
    history: {
      type: Number,
      required: false,
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const name = 'entity_infos_log.delete_after';
    const { updateField } = useModelField(props, emit);

    return {
      name,
      updateField,
    };
  },
};
</script>
