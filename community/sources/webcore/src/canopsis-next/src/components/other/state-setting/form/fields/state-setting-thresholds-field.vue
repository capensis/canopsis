<template lang="pug">
  v-layout(row)
    v-flex(xs4)
      h5.subheading.font-weight-bold {{ $t('stateSetting.criterion') }}: {{ label }}
      state-setting-thresholds-type-field.mt-3(v-field="value.type")
    v-flex
      state-setting-threshold-state-field(
        v-field="value.minor",
        :type="value.type",
        :name="`${name}.minor`",
        :label="$t('stateSetting.states.minor')"
      )
      state-setting-threshold-state-field(
        v-field="value.major",
        :type="value.type",
        :name="`${name}.major`",
        :label="$t('stateSetting.states.major')"
      )
      state-setting-threshold-state-field(
        v-field="value.critical",
        :type="value.type",
        :name="`${name}.critical`",
        :label="$t('stateSetting.states.critical')"
      )
</template>

<script>
import StateSettingThresholdsTypeField from './state-setting-thresholds-type-field.vue';
import StateSettingThresholdStateField from './state-setting-threshold-state-field.vue';

export default {
  inject: ['$validator'],
  components: { StateSettingThresholdStateField, StateSettingThresholdsTypeField },
  model: {
    prop: 'value',
    event: 'input',
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
};
</script>
