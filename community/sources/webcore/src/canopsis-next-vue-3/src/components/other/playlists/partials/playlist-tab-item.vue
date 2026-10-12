<template>
  <v-row class="playlist-tab-item flex-grow-0" no-gutters>
    <v-col class="playlist-tab-item__cell playlist-tab-item__cell--group" cols="4">
      <span class="playlist-tab-item__text">{{ group.title }}</span>
    </v-col>
    <v-col class="playlist-tab-item__cell playlist-tab-item__cell--view" cols="4">
      <span class="playlist-tab-item__text">{{ view.title }}</span>
    </v-col>
    <v-col class="playlist-tab-item__cell playlist-tab-item__cell--tab" cols="4">
      <span class="playlist-tab-item__text">{{ tab.title }}</span>
    </v-col>
  </v-row>
</template>

<script>
export default {
  props: {
    tab: {
      type: Object,
      required: true,
    },
    view: {
      type: Object,
      required: true,
    },
    group: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
.playlist-tab-item {
  width: 100%;
  min-height: 48px;
}

.playlist-tab-item__cell {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 0;
  padding: 12px 28px;
  border: 1px solid var(--v-application-background-darken2);
  color: #fff;
  font-size: 14px;
  line-height: 24px;
  text-align: center;

  &:not(:last-child) {
    margin-right: -1px;
  }

  &--group {
    background-color: #2b3e4f;
  }

  &--view {
    background-color: #425567;
  }

  &--tab {
    background-color: #5a6d80;
  }
}

.playlist-tab-item__text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
