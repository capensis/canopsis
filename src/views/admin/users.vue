<template lang="pug">
  v-container
    div
      div(v-for="user in users", :key="user.id")
        strong {{ user.mail }}
    .fab
      v-speed-dial(
      v-model="fab",
      direction="top",
      transition="slide-y-reverse-transition"
      )
        v-btn(slot="activator", color="green darken-3", dark, fab, v-model="fab")
          v-icon menu
          v-icon close
        v-tooltip(left)
          v-btn(slot="activator", fab, dark, small, color="indigo", @click.stop="showCreateWidgetModal")
            v-icon add
          span Add user
</template>

<script>
import entitiesUserMixin from '@/mixins/entities/user';

export default {
  mixins: [entitiesUserMixin],
  data() {
    return {
      fab: false,
    };
  },
  mounted() {
    this.fetchUsersList();
  },
};
</script>

<style scoped>
  .fab {
    position: fixed;
    bottom: 0;
    right: 0;
  }
</style>
