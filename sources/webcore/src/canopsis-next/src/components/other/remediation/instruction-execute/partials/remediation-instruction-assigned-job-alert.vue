<template lang="pug">
  v-alert.py-1.job-alert(:value="true", type="error", border="left", icon="info")
    v-layout(align-center, justify-space-around)
      span {{ $t('remediationInstructionExecute.jobs.waitAlert') }}
      div.job-alert-buttons
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

    &-buttons {
      display: flex;
    }
  }
</style>
