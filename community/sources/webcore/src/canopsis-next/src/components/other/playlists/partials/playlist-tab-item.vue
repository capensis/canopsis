<template>
  <v-layout
    class="tab-panel-content"
    xs12
  >
    <v-flex
      class="tab-content-block tab-content-block--group white--text"
      xs4
    >
      <span class="tab-content-block__text">{{ group.title }}</span>
    </v-flex>
    <v-flex
      class="tab-content-block tab-content-block--view white--text"
      xs4
    >
      <span class="tab-content-block__text">{{ view.title }}</span>
    </v-flex>
    <v-flex
      class="tab-content-block tab-content-block--tab white--text"
      xs4
    >
      <span class="tab-content-block__text">{{ tab.title }}</span>
    </v-flex>
  </v-layout>
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
.tab-panel-content {
  height: 100%;
}

.tab-content-block {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 0;
  min-height: 48px;
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

  &__text {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
</style>
