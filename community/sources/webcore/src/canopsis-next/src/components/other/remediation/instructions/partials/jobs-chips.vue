<template>
  <v-layout
    class="my-2"
    wrap
  >
    <v-flex
      class="pb-2"
      xs12
    >
      <span class="text-subtitle-1 font-weight-bold">{{ $tc('remediation.instruction.job', 2) }}:</span>
    </v-flex>
    <c-instruction-job-chip
      class="mr-2"
      v-for="job in jobs"
      :key="job._id"
    >
      {{ job.name }}
    </c-instruction-job-chip>
  </v-layout>
</template>

<script>
export default {
  props: {
    jobs: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
