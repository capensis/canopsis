<template>
  <div
    v-if="isOpen"
    class="sidebar-wrapper-overlay"
    @click="closeHandler"
  />

  <v-navigation-drawer
    v-model="isOpen"
    :width="width"
    :style="drawerStyle"
    :scrim="false"
    class="sidebar-wrapper"
    location="right"
    floating
    temporary
  >
    <div v-if="title" class="sidebar-wrapper__sticky-title">
      <v-list :bg-color="sidebar.config.color || 'secondary'">
        <v-list-item>
          <template v-if="sidebar.config.titleIcon" #prepend>
            <v-icon color="white">
              {{ sidebar.config.titleIcon }}
            </v-icon>
          </template>

          <v-list-item-title class="text-h6 text-white">
            {{ title }}
          </v-list-item-title>

          <template #append>
            <v-btn
              class="text-white"
              variant="text"
              icon
              @click.stop="closeHandler"
            >
              <v-icon>close</v-icon>
            </v-btn>
          </template>
        </v-list-item>
      </v-list>

      <v-divider />
    </div>

    <slot />
  </v-navigation-drawer>
</template>

<script>
import { DEFAULT_SIDEBAR_DRAWER_WIDTH } from '@/config';

export default {
  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },

  emits: ['close'],

  data() {
    return {
      ready: false,
    };
  },

  computed: {
    title() {
      return this.sidebar.config.title || (this.sidebar.name ? this.$t(`settings.titles.${this.sidebar.name}`) : '');
    },

    width() {
      return this.sidebar.config.width ?? DEFAULT_SIDEBAR_DRAWER_WIDTH;
    },

    isOpen: {
      get() {
        return Boolean(this.sidebar.name && !this.sidebar.hidden && this.ready);
      },

      set(value) {
        if (!value) {
          this.closeHandler();
        }
      },
    },

    drawerStyle() {
      return {
        top: 0,
        height: '100vh',
        zIndex: 1010,
      };
    },
  },

  mounted() {
    window.requestAnimationFrame(() => {
      this.ready = true;
    });
  },

  methods: {
    closeHandler() {
      this.$emit('close');
    },
  },
};
</script>

<style lang="scss" scoped>
.sidebar-wrapper {
  border-left-width: 0;
  border-inline-start-width: 0;
  overflow: visible;
}

.sidebar-wrapper-overlay {
  position: fixed;
  inset: 0;
  z-index: 1009;
  background: rgb(0, 0, 0);
  opacity: 0.2;
}

.sidebar-wrapper__sticky-title {
  position: sticky;
  top: 0;
  z-index: 2;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
}
</style>
