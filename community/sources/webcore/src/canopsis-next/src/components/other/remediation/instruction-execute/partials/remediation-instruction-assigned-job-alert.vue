<template lang="pug">
  v-alert.py-2.job-alert(:value="true", border="left")
    v-layout(align-center, justify-space-between)
      div.job-alert-content
        v-icon.mr-3(color="white") info
        span {{ $t('remediationInstructionExecute.jobs.waitAlert') }}
      div.job-alert-content
        v-btn.ma-0(
          outline,
          small,
          color="white",
          @click="$emit('skip')"
        ) {{ $t('remediationInstructionExecute.jobs.skip') }}
        v-btn.ma-0.ml-2(
          outline,
          small,
          color="white",
          @click="$emit('await')"
        ) {{ $t('remediationInstructionExecute.jobs.await') }}
</template>

<style lang="scss">
  .job-alert {
    border-top: none;

    &-content {
      display: flex;
      align-items: center;
    }
  }
</style>
