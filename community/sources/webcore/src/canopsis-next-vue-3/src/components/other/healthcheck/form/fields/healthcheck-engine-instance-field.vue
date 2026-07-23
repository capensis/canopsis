<template>
  <c-layout class="healthcheck-engine-instance-field" align="start">
    <v-checkbox
      :model-value="value.enabled"
      :label="label"
      :error-messages="enabledErrors"
      class="healthcheck-engine-instance-field__enabled"
      color="primary"
      hide-details="auto"
      @update:model-value="updateField('enabled', $event)"
    />

    <c-layout class="healthcheck-engine-instance-field__numbers" column>
      <c-layout align="center">
        <div
          :class="{ 'text-disabled': !value.enabled }"
          class="healthcheck-engine-instance-field__number-label"
        >
          {{ $t('common.minimal') }}
        </div>
        <c-number-field
          :value="value.minimal"
          :name="minimalFieldName"
          :label="$t('common.minimal')"
          :disabled="!value.enabled"
          :required="value.enabled"
          :max="+value.optimal"
          :min="0"
          :error-messages="minimalErrors"
          hide-details="auto"
          @update:value="updateField('minimal', $event)"
        />
      </c-layout>

      <c-layout class="mt-2" align="center">
        <div
          :class="{ 'text-disabled': !value.enabled }"
          class="healthcheck-engine-instance-field__number-label"
        >
          {{ $t('common.optimal') }}
        </div>
        <c-number-field
          :value="value.optimal"
          :name="optimalFieldName"
          :label="$t('common.optimal')"
          :disabled="!value.enabled"
          :required="value.enabled"
          :min="+value.minimal"
          :error-messages="optimalErrors"
          hide-details="auto"
          @update:value="updateField('optimal', $event)"
        />
      </c-layout>
    </c-layout>
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
      default: 'engine',
    },
    enabledErrors: {
      type: Array,
      default: () => [],
    },
    minimalErrors: {
      type: Array,
      default: () => [],
    },
    optimalErrors: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['input'],

  computed: {
    enabledFieldName() {
      return `${this.name}.enabled`;
    },

    minimalFieldName() {
      return `${this.name}.minimal`;
    },

    optimalFieldName() {
      return `${this.name}.optimal`;
    },
  },

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
.healthcheck-engine-instance-field {
  gap: 24px;

  &__enabled {
    flex: 0 0 260px;
    max-width: 260px;
  }

  &__numbers {
    max-width: 420px;
    flex: 1 1 auto;
  }

  &__number-label {
    flex: 0 0 120px;
  }
}
</style>
