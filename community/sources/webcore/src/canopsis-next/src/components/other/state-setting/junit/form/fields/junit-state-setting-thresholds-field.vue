<template>
  <v-layout>
    <v-flex xs4>
      <h5 class="text-subtitle-1 font-weight-bold">
        {{ $t('stateSetting.junit.criterion') }}: {{ label }}
      </h5>
      <junit-state-setting-thresholds-type-field
        class="mt-3"
        v-field="value.type"
      />
    </v-flex>
    <v-flex>
      <junit-state-setting-threshold-state-field
        v-field="value.minor"
        :type="value.type"
        :name="`${name}.minor`"
        :label="$t('stateSetting.states.minor')"
      />
      <junit-state-setting-threshold-state-field
        v-field="value.major"
        :type="value.type"
        :name="`${name}.major`"
        :label="$t('stateSetting.states.major')"
      />
      <junit-state-setting-threshold-state-field
        v-field="value.critical"
        :type="value.type"
        :name="`${name}.critical`"
        :label="$t('stateSetting.states.critical')"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import JunitStateSettingThresholdsTypeField from './junit-state-setting-thresholds-type-field.vue';
import JunitStateSettingThresholdStateField from './junit-state-setting-threshold-state-field.vue';

export default {
  components: {
    JunitStateSettingThresholdsTypeField,
    JunitStateSettingThresholdStateField,
  },
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
