<template>
  <c-layout class="remediation-statistics-filters my-4" wrap>
    <c-quick-date-interval-field
      :interval="modelValue.interval"
      :min="minDate"
      class="remediation-statistics-filters__interval mr-4"
      @update:interval="updateField('interval', $event)"
    />
    <c-sampling-field
      :value="modelValue.sampling"
      class="remediation-statistics-filters__sampling mr-4"
      @update:value="updateField('sampling', $event)"
    />
    <remediation-statistics-show-mode-field
      :model-value="modelValue.type"
      class="remediation-statistics-filters__type mr-4"
      @update:model-value="updateField('type', $event)"
    />
    <remediation-statistics-instructions-field
      :model-value="modelValue.instruction"
      class="remediation-statistics-filters__instruction"
      @update:model-value="updateField('instruction', $event)"
    />
  </c-layout>
</template>

<script setup>
import RemediationStatisticsInstructionsField from '../form/fields/remediation-statistics-instructions-field.vue';
import RemediationStatisticsShowModeField from '../form/fields/remediation-statistics-show-mode-field.vue';

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
  minDate: {
    type: Number,
    required: false,
  },
});

const emit = defineEmits(['update:modelValue']);

const updateField = (field, value) => {
  emit('update:modelValue', {
    ...props.modelValue,
    [field]: value,
  });
};
</script>

<style lang="scss" scoped>
.remediation-statistics-filters {
  &__sampling {
    max-width: 200px;
  }

  &__type {
    max-width: 150px;
  }

  &__instruction {
    max-width: 400px;
  }
}
</style>
