<template>
  <c-layout column>
    <v-row>
      <v-col cols="6">
        <junit-state-setting-method-field v-model="method" />
      </v-col>
      <v-col cols="6">
        <state-setting-type-field v-model="type" />
      </v-col>
    </v-row>

    <c-layout v-if="isWorstOfShareMethod" column>
      <c-layout align="center">
        <h4 class="junit-state-setting-form__worst-title text-subtitle-1">
          {{ t('stateSetting.junit.worstLabel') }}
        </h4>
        <c-help-icon
          :text="t('stateSetting.junit.worstHelpText')"
          icon-class="ml-2"
          max-width="220"
          right
        />
      </c-layout>

      <c-layout class="pl-4" gap="4" column>
        <junit-state-setting-thresholds-field
          :value="junitThresholds.skipped"
          :label="t('common.skipped')"
          name="junit_thresholds.skipped"
          @update:value="updateThreshold('skipped', $event)"
        />
        <junit-state-setting-thresholds-field
          :value="junitThresholds.errors"
          :label="t('common.error', 2)"
          name="junit_thresholds.errors"
          @update:value="updateThreshold('errors', $event)"
        />
        <junit-state-setting-thresholds-field
          :value="junitThresholds.failures"
          :label="t('common.failures')"
          name="junit_thresholds.failures"
          @update:value="updateThreshold('failures', $event)"
        />
      </c-layout>
    </c-layout>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { JUNIT_STATE_SETTING_METHODS } from '@/constants';

import { useFormModel } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import StateSettingTypeField from '../../form/fields/state-setting-type-field.vue';

import JunitStateSettingMethodField from './fields/junit-state-setting-method-field.vue';
import JunitStateSettingThresholdsField from './fields/junit-state-setting-thresholds-field.vue';

export default {
  components: {
    StateSettingTypeField,
    JunitStateSettingMethodField,
    JunitStateSettingThresholdsField,
  },

  setup() {
    const { t } = useI18n();
    const method = useFormModel('method');
    const type = useFormModel('type');
    const junitThresholds = useFormModel('junit_thresholds');

    const isWorstOfShareMethod = computed(
      () => method.value === JUNIT_STATE_SETTING_METHODS.worstOfShare,
    );

    const updateThreshold = (field, value) => {
      junitThresholds.value = {
        ...junitThresholds.value,
        [field]: value,
      };
    };

    return {
      t,
      method,
      type,
      junitThresholds,
      isWorstOfShareMethod,
      updateThreshold,
    };
  },
};
</script>

<style lang="scss" scoped>
.junit-state-setting-form {
  &__worst-title {
    font-weight: 700;
  }
}
</style>
