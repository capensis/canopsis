<template>
  <c-radio-group
    :model-value="value"
    hide-details
    inline
    @update:model-value="$emit('update:value', $event)"
  >
    <v-radio
      :label="$t('view.duplicateAsPrivate')"
      :value="true"
      color="primary"
    />
    <v-radio
      :label="$t('view.duplicateAsRegular')"
      :value="false"
      color="primary"
    />
  </c-radio-group>
</template>

<script>
export default {
  props: {
    value: {
      type: Boolean,
      required: true,
    },
  },

  emits: ['update:value'],
};
</script>
