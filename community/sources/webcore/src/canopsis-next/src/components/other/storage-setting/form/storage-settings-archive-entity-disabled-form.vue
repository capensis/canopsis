<template>
  <c-checkbox
    :model-value="modelValue.with_dependencies"
    color="primary"
    @update:model-value="updateField('with_dependencies', $event)"
  >
    <template #label="">
      {{ $t('storageSetting.entity.archiveDependencies') }}
      <c-help-icon
        :text="$t('storageSetting.entity.archiveDependenciesHelp')"
        icon-class="ml-2"
        color="info"
        max-width="300"
        top
      />
    </template>
  </c-checkbox>
</template>

<script>
import { useModelField } from '@/hooks/form';

export default {
  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    return { updateField };
  },
};
</script>
