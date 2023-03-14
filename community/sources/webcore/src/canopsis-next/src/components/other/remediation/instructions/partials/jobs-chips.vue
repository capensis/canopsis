<template lang="pug">
  v-layout.my-2(wrap)
    v-flex.pb-2(xs12)
      span.subheading.font-weight-bold {{ $tc('remediation.instruction.job', 2) }}:
    c-instruction-job-chip.mr-2(v-for="job in jobs", :key="job._id") {{ job.name }}
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
