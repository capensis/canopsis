<template>
  <v-menu
    :left="left"
    :top="top"
    :offset-x="left"
    :offset-y="top"
    :nudge-left="left ? 10 : 0"
    :nudge-top="top ? 10 : 0"
    class="view-screen-mode-btn"
  >
    <template #activator="{ on }">
      <v-btn
        :small="small"
        class="view-fullscreen-btn__button"
        fab
        dark
        v-on="on"
      >
        <v-icon :small="small">
          {{ icon }}
        </v-icon>
      </v-btn>
    </template>
    <v-list class="view-screen-mode-btn__list grey darken-3 pa-0" dark>
      <v-tooltip
        v-for="item in availableScreenModes"
        :key="`tooltip-${item.key}`"
        left
      >
        <template #activator="{ on: tooltipOn }">
          <v-list-item
            v-bind="item.bind"
            active-class="view-screen-mode-btn__item--active"
            v-on="{ ...item.on, ...tooltipOn }"
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
import { computed } from 'vue';

import { VIEW_SCREEN_MODES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    value: {
      type: String,
      default: VIEW_SCREEN_MODES.default,
    },
    small: {
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
  setup(props, { emit }) {
    const { t } = useI18n();

    const changeMode = (newMode) => {
      if (newMode === props.value) {
        return;
      }

      emit('input', newMode);
    };

    const availableScreenModes = computed(() => Object.values(VIEW_SCREEN_MODES).map(mode => ({
      key: mode,
      title: t(`view.screenMode.${mode}.title`),
      tooltip: t(`view.screenMode.${mode}.tooltip`),
      bind: {
        inputValue: props.value === mode,
      },
      on: {
        click: () => changeMode(mode),
      },
    })));

    const icon = computed(() => ([
      VIEW_SCREEN_MODES.fullscreen,
      VIEW_SCREEN_MODES.kioskFullscreen,
    ].includes(props.value)
      ? 'fullscreen_exit'
      : 'fullscreen'));

    return {
      availableScreenModes,
      icon,
    };
  },
};
</script>

<style lang="scss">
.view-screen-mode-btn {
  &__button  {
    border-color: #212121 !important;
    background-color: #212121 !important;

    .theme--dark & {
      border-color: #979797 !important;
      background-color: #979797 !important;
    }
  }

  &__list {
    .view-screen-mode-btn__item--active {
      background-color: var(--v-primary-darken1);
    }
  }
}
</style>
