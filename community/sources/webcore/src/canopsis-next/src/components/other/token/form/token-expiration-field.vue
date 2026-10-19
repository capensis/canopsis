<template>
  <div class="token-expiration-field">
    <c-form-block-row :label="label" width="25%">
      <c-enabled-field
        :model-value="enabled"
        :label="label"
        inline
        @update:model-value="$emit('update:enabled', $event)"
      />

      <c-duration-field
        v-if="enabled"
        :duration="duration"
        class="token-expiration-field__duration"
        name="duration"
        long
        @update:duration="$emit('update:duration', $event)"
      />
    </c-form-block-row>
  </div>
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      required: true,
    },
    enabled: {
      type: Boolean,
      required: true,
    },
    duration: {
      type: Object,
      required: true,
    },
  },

  emits: ['update:enabled', 'update:duration'],
};
</script>

<style lang="scss" scoped>
.token-expiration-field {
  overflow: hidden;
  border: 1px solid var(--v-application-background-darken2);
  border-radius: 5px;

  &__duration {
    max-width: 420px;
    margin-top: 12px;
  }

  :deep(.c-form-block-row__label) {
    align-items: center;
    padding: 12px 16px;
  }
}
</style>
