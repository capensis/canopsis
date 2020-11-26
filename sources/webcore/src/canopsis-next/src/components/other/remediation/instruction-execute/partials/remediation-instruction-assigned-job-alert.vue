<template lang="pug">
  v-alert.py-1.job-alert(:value="true", type="error", border="left", icon="info")
    v-layout(align-center, justify-space-around)
      v-flex(xs8)
        span {{ $t('remediationInstructionExecute.jobs.waitAlert') }}
      v-flex(xs4)
        v-layout(row)
          v-btn(
            outline,
            small,
            block,
            color="white",
            @click="$emit('skip')"
          ) {{ $t('remediationInstructionExecute.jobs.skip') }}
          v-btn.ml-2(
            outline,
            small,
            block,
            color="white",
            @click="$emit('await')"
          ) {{ $t('remediationInstructionExecute.jobs.await') }}
</template>

<style lang="scss">
  .job-alert {
    border-top: none;
  }
</style>
