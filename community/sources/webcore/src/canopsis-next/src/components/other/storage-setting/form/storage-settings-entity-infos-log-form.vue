<template>
  <c-information-block :title="$t('storageSetting.entityInfosLog.title')">
    <c-enabled-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.entityInfosLog.deleteAfter')"
      :name="name"
    />
  </c-information-block>
</template>

<script>
export default {
  model: {
    prop: 'form',
    event: 'input',
  },
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
  setup() {
    const name = 'entity_infos_log.delete_after';

    return {
      name,
    };
  },
};
</script>
