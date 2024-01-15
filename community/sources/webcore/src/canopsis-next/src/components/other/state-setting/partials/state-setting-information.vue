<template>
  <v-layout
    class="py-4"
    column
  >
    <state-setting-information-row
      :label="$t('common.name')"
      :value="stateSetting.type"
    />
    <state-setting-information-row
      :label="$t('common.method')"
      :value="$t(`stateSetting.methods.${stateSetting.method}`)"
    />
  </v-layout>
</template>

<script>
import StateSettingInformationRow from './state-setting-information-row.vue';

export default {
  components: { StateSettingInformationRow },
  props: {
    stateSetting: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
