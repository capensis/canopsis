<template>
  <c-layout class="healthcheck-enabled-limit-field" align="center">
    <v-checkbox
      :model-value="value.enabled"
      :label="label"
      :error-messages="enabledErrors"
      class="healthcheck-enabled-limit-field__enabled"
      color="primary"
      hide-details="auto"
      @update:model-value="updateField('enabled', $event)"
    />
    <c-number-field
      :value="value.limit"
      :name="name"
      :label="label"
      :disabled="!value.enabled"
      :required="value.enabled"
      :min="min"
      :error-messages="limitErrors"
      class="healthcheck-enabled-limit-field__limit"
      hide-details="auto"
      @update:value="updateField('limit', $event)"
    />
  </c-layout>
</template>

<script>
export default {
  props: {
    value: {
      type: Object,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'limit',
    },
    min: {
      type: Number,
      default: 1,
    },
    enabledErrors: {
      type: Array,
      default: () => [],
    },
    limitErrors: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['input'],

  methods: {
    updateField(field, value) {
      this.$emit('input', {
        ...this.value,
        [field]: value,
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.healthcheck-enabled-limit-field {
  gap: 24px;

  &__enabled {
    flex: 0 0 260px;
    max-width: 260px;
  }

  &__limit {
    max-width: 360px;
  }
}
</style>
