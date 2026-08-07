<template>
  <div
    :class="{ 'view-fab-btns__fullscreen': isFullscreen }"
    class="view-fab-btns fab ma-2"
  >
    <c-layout gap="3">
      <view-scroll-top-btn />
      <view-periodic-refresh-btn />
      <view-screen-mode-btn
        v-if="isFullscreen"
        :value="activeViewScreenMode"
        large
        top
        @input="changeScreenMode"
      />
      <v-menu
        v-else
        v-model="isMenuOpen"
        :close-on-content-click="false"
        location="top center"
        transition="slide-y-reverse-transition"
      >
        <template #activator="{ props: activatorProps }">
          <v-btn
            v-bind="activatorProps"
            color="primary"
            size="x-large"
            icon
          >
            <v-icon :icon="isMenuOpen ? 'close' : 'menu'" />
          </v-btn>
        </template>

        <c-layout
          class="view-fab-btns__list"
          align="center"
          gap="2"
          column
        >
          <v-tooltip
            v-if="updatable"
            location="left"
          >
            <template #activator="{ props: tooltipProps }">
              <v-btn
                v-bind="tooltipProps"
                color="green"
                icon
                @click.stop="showCreateTabModal"
              >
                <v-icon
                  icon="add"
                  size="small"
                />
              </v-btn>
            </template>
            <span>{{ $t('common.addTab') }}</span>
          </v-tooltip>

          <v-tooltip
            v-if="updatable"
            location="left"
          >
            <template #activator="{ props: tooltipProps }">
              <v-btn
                v-bind="tooltipProps"
                class="view-fab-btns__add-widget-btn"
                icon
                @click.stop="showCreateWidgetModal"
              >
                <v-icon
                  icon="add"
                  size="small"
                />
              </v-btn>
            </template>
            <span>{{ $t('common.addWidget') }}</span>
          </v-tooltip>

          <view-editing-btn v-if="updatable" />

          <view-screen-mode-btn
            :value="activeViewScreenMode"
            left
            @input="changeScreenMode"
          />

          <view-share-link-btn
            v-if="sharable"
            :view="view"
            :tab="activeTab"
          />
        </c-layout>
      </v-menu>
    </c-layout>
  </div>
</template>

<script>
import { onBeforeUnmount, onMounted, ref } from 'vue';

import {
  FULLSCREEN_MODES_TO_DEFAULT_SCREEN_MODES,
  KEYS_TO_VIEW_SCREEN_MODES,
  MODALS,
  VIEW_SCREEN_MODES,
} from '@/constants';

import { useActiveView } from '@/hooks/store/modules/active-view';
import { useFullscreen } from '@/hooks/fullscreen';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useView } from '@/hooks/store/modules/view';
import { useViewRouter } from '@/hooks/view/router';

import ViewEditingBtn from './view-editing-btn.vue';
import ViewPeriodicRefreshBtn from './view-periodic-refresh-btn.vue';
import ViewScreenModeBtn from './view-screen-mode-btn.vue';
import ViewShareLinkBtn from './view-share-link-btn.vue';
import ViewScrollTopBtn from './view-scroll-top-btn.vue';

export default {
  components: {
    ViewEditingBtn,
    ViewPeriodicRefreshBtn,
    ViewScreenModeBtn,
    ViewShareLinkBtn,
    ViewScrollTopBtn,
  },
  props: {
    activeTab: {
      type: Object,
      required: false,
      default: null,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    sharable: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const isFullscreen = ref(false);
    const isMenuOpen = ref(false);

    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const fullscreen = useFullscreen();
    const { route, redirectToFirstTab } = useViewRouter();
    const { createViewTab } = useView();
    const {
      view,
      activeViewScreenMode,
      fetchActiveView,
      toggleEditing,
      setActiveViewScreenMode,
    } = useActiveView();

    const enterFullscreen = () => {
      if (!props.activeTab) {
        popups.warning({ text: t('view.errors.emptyTabs') });
        return;
      }

      const element = document.querySelector('[data-app]') ?? document.getElementById('app');
      const viewElement = document.getElementById(`view-tab-${props.activeTab._id}`);

      if (!element || !viewElement) {
        return;
      }

      fullscreen.request(element, {
        fullscreenClass: 'full-screen',
        callback: (value) => {
          viewElement.classList.toggle('view-fullscreen', value);
          isFullscreen.value = value;

          const newMode = FULLSCREEN_MODES_TO_DEFAULT_SCREEN_MODES[activeViewScreenMode.value];

          if (!value && newMode) {
            setActiveViewScreenMode(newMode);
          }
        },
      });
    };

    const exitFullscreen = () => fullscreen.exit();

    const changeScreenMode = (newMode) => {
      setActiveViewScreenMode(newMode);

      switch (newMode) {
        case VIEW_SCREEN_MODES.fullscreen:
        case VIEW_SCREEN_MODES.kioskFullscreen:
          enterFullscreen();
          return;

        default:
          exitFullscreen();
      }
    };

    const keyDownListener = (event) => {
      if (event.code === 'KeyE' && (event.ctrlKey || event.metaKey) && props.updatable) {
        toggleEditing();
        event.preventDefault();
        return;
      }

      if ((event.altKey || event.metaKey) && event.shiftKey) {
        const newMode = KEYS_TO_VIEW_SCREEN_MODES[event.code];

        if (newMode) {
          changeScreenMode(newMode);
          event.preventDefault();
        }
      }
    };

    const showCreateWidgetModal = () => {
      if (!props.activeTab) {
        popups.warning({ text: t('view.errors.emptyTabs') });
        return;
      }

      modals.show({
        name: MODALS.createWidget,
        config: {
          tab: props.activeTab,
        },
      });
    };

    const showCreateTabModal = () => {
      modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: t('modals.viewTab.create.title'),
          field: {
            name: 'text',
            label: t('modals.viewTab.fields.title'),
            validationRules: {
              required: true,
            },
          },
          action: async (title) => {
            await createViewTab({
              data: {
                view: view.value._id,
                title,
              },
            });
            await fetchActiveView();

            if (!route.query.tabId) {
              await redirectToFirstTab();
            }
          },
        },
      });
    };

    onMounted(() => document.addEventListener('keydown', keyDownListener));

    onBeforeUnmount(() => {
      document.removeEventListener('keydown', keyDownListener);
      exitFullscreen();
    });

    return {
      isFullscreen,
      isMenuOpen,
      view,
      activeViewScreenMode,
      changeScreenMode,
      showCreateWidgetModal,
      showCreateTabModal,
    };
  },
};
</script>

<style lang="scss">
.view-fab-btns {
  &__list {
    align-items: center;
    padding: 8px 0 8px 8px;
  }

  &__add-widget-btn {
    border-color: #3f51b5 !important;
    background-color: #3f51b5 !important;
    color: #fff !important;

    .v-theme--dark & {
      border-color: #2196f3 !important;
      background-color: #2196f3 !important;
    }
  }

  &__fullscreen {
    z-index: 8;
  }
}

.view-fullscreen {
  overflow: auto;
  position: fixed;
  inset: 0;
  z-index: 7;
  background: rgb(var(--v-theme-surface));
}
</style>
