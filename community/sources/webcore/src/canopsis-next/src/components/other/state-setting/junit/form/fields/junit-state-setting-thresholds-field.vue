<template>
  <c-layout class="junit-state-setting-thresholds-field" gap="6">
    <c-layout class="junit-state-setting-thresholds-field__criterion" column>
      <h5 class="junit-state-setting-thresholds-field__title text-subtitle-1">
        {{ t('stateSetting.junit.criterion') }}: {{ label }}
      </h5>
      <junit-state-setting-thresholds-type-field
        :model-value="value.type"
        :name="`${name}.type`"
        class="mt-2"
        @update:model-value="updateValueField('type', $event)"
      />
    </c-layout>

    <c-layout class="junit-state-setting-thresholds-field__states" gap="1" column>
      <junit-state-setting-threshold-state-field
        :model-value="value.minor"
        :type="value.type"
        :max="value.major"
        :name="`${name}.minor`"
        :label="t('common.stateTypes.minor')"
        @update:model-value="updateValueField('minor', $event)"
      />
      <junit-state-setting-threshold-state-field
        :model-value="value.major"
        :type="value.type"
        :max="value.critical"
        :name="`${name}.major`"
        :label="t('common.stateTypes.major')"
        @update:model-value="updateValueField('major', $event)"
      />
      <junit-state-setting-threshold-state-field
        :model-value="value.critical"
        :type="value.type"
        :name="`${name}.critical`"
        :label="t('common.stateTypes.critical')"
        @update:model-value="updateValueField('critical', $event)"
      />
    </c-layout>
  </c-layout>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

import JunitStateSettingThresholdStateField from './junit-state-setting-threshold-state-field.vue';
import JunitStateSettingThresholdsTypeField from './junit-state-setting-thresholds-type-field.vue';

export default {
  components: {
    JunitStateSettingThresholdStateField,
    JunitStateSettingThresholdsTypeField,
  },

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
      default: '',
    },
  },

  emits: ['update:value'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const updateValueField = (field, value) => {
      emit('update:value', {
        ...props.value,
        [field]: value,
      });
    };

    return {
      t,
      updateValueField,
    };
  },
};
</script>

<style lang="scss" scoped>
.junit-state-setting-thresholds-field {
  width: 100%;

  &__criterion,
  &__states {
    flex: 1 1 0;
    min-width: 0;
    align-self: start;
  }

  &__title {
    font-weight: 700;
  }
}
</style>
