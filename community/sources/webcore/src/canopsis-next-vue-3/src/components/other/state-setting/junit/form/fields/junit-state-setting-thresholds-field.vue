<template>
  <v-row>
    <v-col cols="4">
      <h5 class="junit-state-setting-thresholds-field__title text-subtitle-1">
        {{ t('stateSetting.junit.criterion') }}: {{ label }}
      </h5>
      <junit-state-setting-thresholds-type-field
        :model-value="value.type"
        :name="`${name}.type`"
        class="mt-3"
        @update:model-value="updateValueField('type', $event)"
      />
    </v-col>

    <v-col>
      <junit-state-setting-threshold-state-field
        :model-value="value.minor"
        :type="value.type"
        :name="`${name}.minor`"
        :label="t('common.stateTypes.minor')"
        @update:model-value="updateValueField('minor', $event)"
      />
      <junit-state-setting-threshold-state-field
        :model-value="value.major"
        :type="value.type"
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
    </v-col>
  </v-row>
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
.junit-state-setting-thresholds-field__title {
  font-weight: 700;
}
</style>
