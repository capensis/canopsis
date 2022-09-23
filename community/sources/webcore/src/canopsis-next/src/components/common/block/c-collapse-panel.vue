<template lang="pug">
  v-expansion-panel.c-collapse-panel
    v-expansion-panel-content(:class="{ error }", :style="{ backgroundColor: color }", :lazy="lazy")
      template(#actions="")
        v-icon(color="white") {{ icon }}
      template(#header="")
        slot(name="header")
      div.white
        slot
</template>

<script>
export default {
  props: {
    color: {
      type: String,
      default: '',
    },
    icon: {
      type: String,
      default: '$vuetify.icons.expand',
    },
    lazy: {
      type: Boolean,
      default: false,
    },
    error: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss">
.c-collapse-panel {
  border-radius: 5px;
  overflow: hidden;
}
</style>
