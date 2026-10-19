<template>
  <div :class="{ 'view-fab-btns__fullscreen': isFullscreen }" class="view-fab-btns fab ma-2">
    <c-layout gap="3">
      <view-scroll-top-btn />
      <view-executions-btn v-if="hasAccessToExecuteInstruction" />
      <view-periodic-refresh-btn />
      <view-screen-mode-btn
        v-if="isFullscreen"
        :value="activeViewScreenMode"
        large
        top
        @input="changeScreenMode"
      />
      <c-speed-dial
        v-else
        :close-on-content-click="false"
        location="top center"
      >
        <template #activator="{ props: activatorProps, isActive }">
          <v-btn
            v-bind="activatorProps"
            :class="{ 'c-speed-dial__btn--active': isActive }"
            class="c-speed-dial__btn"
            color="primary"
            size="x-large"
            icon
          >
            <v-icon icon="menu" />
            <v-icon icon="close" />
          </v-btn>
        </template>

        <div v-if="sharable" key="share">
          <view-share-link-btn :view="view" :tab="activeTab" />
        </div>

        <div key="screen-mode">
          <view-screen-mode-btn :value="activeViewScreenMode" left @input="changeScreenMode" />
        </div>

        <div v-if="updatable" key="editing">
          <view-editing-btn />
        </div>

        <div v-if="updatable" key="add-widget">
          <c-tooltip location="left">
            <template #activator="{ props: tooltipProps }">
              <v-btn
                v-bind="tooltipProps"
                class="view-fab-btns__add-widget-btn"
                size="default"
                icon
                @click.stop="showCreateWidgetModal"
              >
                <v-icon icon="add" size="small" />
              </v-btn>
            </template>
            <span>{{ $t('common.addWidget') }}</span>
          </c-tooltip>
        </div>

        <div v-if="updatable" key="add-tab">
          <c-tooltip location="left">
            <template #activator="{ props: tooltipProps }">
              <v-btn
                v-bind="tooltipProps"
                color="green"
                size="default"
                icon
                @click.stop="showCreateTabModal"
              >
                <v-icon icon="add" size="small" />
              </v-btn>
            </template>
            <span>{{ $t('common.addTab') }}</span>
          </c-tooltip>
        </div>
      </c-speed-dial>
    </c-layout>
  </div>
</template>

<script>
import { onBeforeUnmount, onMounted, ref } from 'vue';

import {
  FULLSCREEN_MODES_TO_DEFAULT_SCREEN_MODES,
  KEYS_TO_VIEW_SCREEN_MODES,
  MODALS,
  ROUTES_NAMES,
  USER_PERMISSIONS,
  VIEW_SCREEN_MODES,
} from '@/constants';

import { useCanPermission } from '@/hooks/auth';
import { useFullscreen } from '@/hooks/fullscreen';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useActiveView } from '@/hooks/store/modules/active-view';
import { useInfo } from '@/hooks/store/modules/info';
import { useView } from '@/hooks/store/modules/view';
import { useViewRouter } from '@/hooks/view/router';

import ViewEditingBtn from './view-editing-btn.vue';
import ViewExecutionsBtn from './view-executions-btn.vue';
import ViewPeriodicRefreshBtn from './view-periodic-refresh-btn.vue';
import ViewScreenModeBtn from './view-screen-mode-btn.vue';
import ViewScrollTopBtn from './view-scroll-top-btn.vue';
import ViewShareLinkBtn from './view-share-link-btn.vue';

export default {
  components: {
    ViewEditingBtn,
    ViewExecutionsBtn,
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

    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const fullscreen = useFullscreen();
    const { route, redirectToFirstTab } = useViewRouter();
    const { createViewTab } = useView();
    const { showHeaderOnKioskMode } = useInfo();
    const { view, activeViewScreenMode, fetchActiveView, toggleEditing, setActiveViewScreenMode } =
      useActiveView();
    const { hasAccess: hasAccessToExecuteInstruction } = useCanPermission(
      USER_PERMISSIONS.business.alarmsList.actions.executeInstruction,
    );

    const isHeaderShown = () => {
      if (
        [VIEW_SCREEN_MODES.fullscreen, VIEW_SCREEN_MODES.kioskFullscreen].includes(
          activeViewScreenMode.value,
        )
      ) {
        return false;
      }

      if (
        [VIEW_SCREEN_MODES.kiosk, VIEW_SCREEN_MODES.kioskFullscreen].includes(
          activeViewScreenMode.value,
        ) ||
        route.name === ROUTES_NAMES.viewKiosk
      ) {
        return showHeaderOnKioskMode.value;
      }

      return !route.meta?.hideHeader;
    };

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
          viewElement.classList.toggle('view-fullscreen--with-top-bar', value && isHeaderShown());
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
      view,
      activeViewScreenMode,
      hasAccessToExecuteInstruction,
      changeScreenMode,
      showCreateWidgetModal,
      showCreateTabModal,
    };
  },
};
</script>

<style lang="scss">
.view-fab-btns {
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

  &--with-top-bar {
    top: 48px;
  }
}
</style>
