<template>
  <v-layout>
    <v-subheader class="primary--text">
      Instruction name
    </v-subheader>
    <span>Alarm name</span>
    <v-layout class="gap-2">
      <strong>Jobs:</strong>
      <active-remediation-executions-item-jobs :jobs="jobs" />
    </v-layout>
    <span>
      <strong>Step 1a:&nbsp;</strong>
      <span>Step name</span>
    </span>
  </v-layout>
</template>

<script>
export default {
  setup() {
    return {};
  },
};
</script>
