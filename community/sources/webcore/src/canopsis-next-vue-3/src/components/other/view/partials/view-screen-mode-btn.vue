<template>
  <v-menu
    :location="location"
    :offset="10"
    class="view-screen-mode-btn"
  >
    <template #activator="{ props: activatorProps }">
      <v-btn
        v-bind="activatorProps"
        :size="large ? 'x-large' : undefined"
        class="view-screen-mode-btn__button"
        icon
      >
        <v-icon
          :icon="icon"
          :size="large ? 'default' : 'small'"
        />
      </v-btn>
    </template>

    <v-list class="view-screen-mode-btn__list bg-grey-darken-3 pa-0">
      <v-tooltip
        v-for="item in availableScreenModes"
        :key="`tooltip-${item.key}`"
        location="left"
      >
        <template #activator="{ props }">
          <v-list-item
            v-bind="props"
            :active="item.active"
            active-class="view-screen-mode-btn__item--active"
            @click="item.click"
          >
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item>
        </template>
        <span>{{ item.tooltip }}</span>
      </v-tooltip>
    </v-list>
  </v-menu>
</template>

<script>
import { VIEW_SCREEN_MODES } from '@/constants';

export default {
  props: {
    value: {
      type: String,
      default: VIEW_SCREEN_MODES.default,
    },
    large: {
      type: Boolean,
      default: false,
    },
    left: {
      type: Boolean,
      default: false,
    },
    top: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['input', 'update:value'],

  computed: {
    location() {
      if (this.left) {
        return 'left center';
      }

      return this.top ? 'top center' : 'bottom center';
    },

    availableScreenModes() {
      return Object.values(VIEW_SCREEN_MODES).map(mode => ({
        key: mode,
        title: this.$t(`view.screenMode.${mode}.title`),
        tooltip: this.$t(`view.screenMode.${mode}.tooltip`),
        active: this.value === mode,
        click: () => this.changeMode(mode),
      }));
    },

    icon() {
      return [
        VIEW_SCREEN_MODES.fullscreen,
        VIEW_SCREEN_MODES.kioskFullscreen,
      ].includes(this.value)
        ? 'fullscreen_exit'
        : 'fullscreen';
    },
  },

  methods: {
    changeMode(newMode) {
      if (newMode === this.value) {
        return;
      }

      this.$emit('input', newMode);
      this.$emit('update:value', newMode);
    },
  },
};
</script>

<style lang="scss" scoped>
.view-screen-mode-btn__button {
  border-color: #212121 !important;
  background-color: #212121 !important;
  color: #fff !important;

  .v-theme--dark & {
    border-color: #979797 !important;
    background-color: #979797 !important;
  }
}

.view-screen-mode-btn__list {
  color: #fff;

  :deep(.view-screen-mode-btn__item--active) {
    background-color: rgb(var(--v-theme-primary));
  }
}
</style>
