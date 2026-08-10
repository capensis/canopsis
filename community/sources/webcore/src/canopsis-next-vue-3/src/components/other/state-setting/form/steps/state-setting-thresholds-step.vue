<template>
  <c-layout class="state-setting-thresholds-step" column>
    <state-setting-threshold-field
      v-for="key in sortedThresholdsKeys"
      :key="key"
      :condition="thresholds[key]"
      :label="getLabel(key)"
      :name="`${name}.${key}`"
      :state="key"
      @update:condition="updateThreshold(key, $event)"
    />

    <c-alert v-if="hasError" type="error">
      {{ t('stateSetting.conditionsError') }}
    </c-alert>
  </c-layout>
</template>

<script>
import { computed, onBeforeUnmount, onMounted, watch } from 'vue';

import { ALARM_STATES } from '@/constants';
import { useFormModel } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useValidationAttachRequired } from '@/hooks/validator/validation-attach-required';
import { useValidationHeader } from '@/hooks/validator/validation-header';

import StateSettingThresholdField from '../fields/state-setting-threshold-field.vue';

export default {
  components: { StateSettingThresholdField },

  props: {
    name: {
      type: String,
      default: 'state_thresholds',
    },
  },

  setup(props) {
    const { t } = useI18n();
    const thresholds = useFormModel(props.name);
    const { hasAnyError } = useValidationHeader();
    const { validator, attachRequiredRule, detachRequiredRule, validateRequiredRule } =
      useValidationAttachRequired(props.name);

    const sortedThresholdsKeys = computed(() => Object.keys(thresholds.value ?? {}).reverse());
    const hasError = computed(() => validator?.errors?.has(props.name) ?? false);

    const getLabel = key => t(`common.stateTypes.${ALARM_STATES[key]}`);

    const requiredRuleGetter = () =>
      Object.values(thresholds.value ?? {}).some(({ enabled }) => enabled);

    const updateThreshold = (key, value) => {
      thresholds.value = {
        ...thresholds.value,
        [key]: value,
      };

      validator?.errors?.clear(props.name);
    };

    watch(thresholds, () => validateRequiredRule(), { deep: true });

    onMounted(() => attachRequiredRule(requiredRuleGetter));
    onBeforeUnmount(() => detachRequiredRule());

    return {
      t,
      thresholds,
      sortedThresholdsKeys,
      hasAnyError,
      hasError,
      getLabel,
      updateThreshold,
    };
  },
};
</script>

<style lang="scss" scoped>
.state-setting-thresholds-step {
  gap: 16px;
}
</style>
