<template lang="pug">
    v-toolbar(
      dense,
      fixed,
      clipped-left,
      app,
      color="blue darken-4",
    )
      v-toolbar-side-icon.white--text(@click="toggleSideBar")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('app');

/**
 * Component for the top bar of the application
 */
export default {
  methods: {
    ...mapActions([
      'toggleSideBar',
    ]),
  },
};
</script>
