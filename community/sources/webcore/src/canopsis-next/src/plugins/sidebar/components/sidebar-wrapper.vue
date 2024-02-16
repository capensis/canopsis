<template>
  <v-navigation-drawer
    v-model="isOpen"
    :ignore-click-outside="hasMaximizedModal"
    :custom-close-conditional="closeCondition"
    :width="450"
    right
    fixed
    temporary
  >
    <div v-if="title">
      <v-list color="secondary">
        <v-list-item>
          <v-list-item-title class="white--text text-subtitle-1">
            {{ title }}
          </v-list-item-title>
          <v-btn
            icon
            @click.stop="closeHandler"
          >
            <v-icon color="white">
              close
            </v-icon>
          </v-btn>
        </v-list-item>
      </v-list>
      <v-divider />
    </div>
    <slot />
  </v-navigation-drawer>
</template>

<script>
/**
 * Wrapper for each sidebar
 */
export default {
  inject: ['$clickOutside'],
  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      ready: false,
    };
  },
  computed: {
    hasMaximizedModal() {
      return this.$store.getters[`${this.$modals.moduleName}/hasMaximizedModal`];
    },

    title() {
      return this.sidebar.name ? this.$t(`settings.titles.${this.sidebar.name}`) : '';
    },

    isOpen: {
      get() {
        return this.sidebar.name && !this.sidebar.hidden && this.ready;
      },
      set(value) {
        if (!value) {
          this.$sidebar.hide();
        }
      },
    },
  },
  mounted() {
    this.ready = true;
  },
  methods: {
    closeHandler() {
      if (this.closeCondition()) {
        this.$sidebar.hide();
      }
    },

    closeCondition(...args) {
      return this.$clickOutside.call(...args);
    },
  },
};
</script>
