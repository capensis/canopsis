<template lang="pug">
  v-tooltip(left)
    v-btn(
      slot="activator",
      color="secondary",
      data-test="refreshButton",
      fab,
      dark,
      icon,
      @click.stop="$emit('click', $event)"
    )
      v-icon refresh
    span {{ $t('common.refresh') }}
</template>
