<template lang="pug">
  v-tooltip(left)
    v-btn(slot="activator", fab, dark, icon, color="secondary", @click.stop="$emit('click', $event)")
      v-icon refresh
    span {{ $t('common.refresh') }}
</template>
