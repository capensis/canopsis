<template lang="pug">
    v-toolbar(
      dense,
      fixed,
      clipped-left,
      app,
      color="blue darken-4",
    )
      v-toolbar-side-icon(
        class="white--text",
        @click="toggleSideBar",
      )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('app');

export default {
  name: 'TopBar',
  methods: {
    ...mapActions([
      'toggleSideBar',
    ]),
  },
};
</script>
