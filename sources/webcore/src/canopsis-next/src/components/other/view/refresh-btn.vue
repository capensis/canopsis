<template lang="pug">
  v-tooltip(left)
    v-btn(slot="activator", fab, dark, icon, color="secondary", @click.stop="action")
      v-icon refresh
    span {{ $t('common.refresh') }}
</template>

<script>
export default {
  props: {
    action: {
      type: Function,
      required: true,
    },
  },
};
</script>
