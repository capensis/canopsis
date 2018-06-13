<template lang="pug">
  v-layout(align-center, wrap)
    div Selection :
    v-tooltip(bottom)
      v-btn(slot="activator", small, icon)
        v-icon pause
      span Add a PBehaviour
    v-tooltip(bottom)
      v-btn(slot="activator", small, icon)
        v-icon check
      span Fast ack
    v-tooltip(bottom)
      v-btn(slot="activator", small, icon)
        v-icon playlist_add_check
      span Ack
    v-tooltip(bottom)
      v-btn(slot="activator", small, icon)
        v-icon block
      span Cancel Ack
    v-tooltip(bottom)
      v-btn(slot="activator", small, icon)
        v-icon undo
      span Restore alarm
</template>
