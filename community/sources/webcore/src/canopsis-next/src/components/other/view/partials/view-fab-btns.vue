<template>
  <div
    :class="{ 'view-fab-btns__fullscreen': isFullscreen }"
    class="view-fab-btns fab ma-2"
  >
    <v-layout>
      <v-flex class="mr-3">
        <view-scroll-top-btn />
      </v-flex>
      <v-flex class="mr-3">
        <view-periodic-refresh-btn />
      </v-flex>
      <v-flex v-if="isFullscreen">
        <view-screen-mode-btn
          :value="activeViewScreenMode"
          top
          @input="changeScreenMode"
        />
      </v-flex>
      <c-speed-dial
        v-else
        direction="top"
        transition="slide-y-reverse-transition"
        eager
      >
        <template #activator="{ bind }">
          <v-btn
            v-bind="bind"
            color="primary"
            dark
            fab
          >
            <v-icon>menu</v-icon>
            <v-icon>close</v-icon>
          </v-btn>
        </template>
        <view-share-link-btn
          v-if="sharable"
          :view="view"
          :tab="activeTab"
        />
        <view-screen-mode-btn
          :value="activeViewScreenMode"
          left
          small
          @input="changeScreenMode"
        />
        <view-editing-btn v-if="updatable" />
        <v-tooltip left>
          <template #activator="{ on }">
            <v-btn
              v-if="updatable"
              class="view-fab-btns__add-widget-btn"
              fab
              dark
              small
              v-on="on"
              @click.stop="showCreateWidgetModal"
            >
              <v-icon small>
                add
              </v-icon>
            </v-btn>
          </template>
          <span>{{ $t('common.addWidget') }}</span>
        </v-tooltip>
        <v-tooltip left>
          <template #activator="{ on }">
            <v-btn
              v-if="updatable"
              color="green"
              fab
              dark
              small
              v-on="on"
              @click.stop="showCreateTabModal"
            >
              <v-icon small>
                add
              </v-icon>
            </v-btn>
          </template>
          <span>{{ $t('common.addTab') }}</span>
        </v-tooltip>
      </c-speed-dial>
    </v-layout>
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount } from 'vue';

import {
  MODALS,
  KEYS_TO_VIEW_SCREEN_MODES,
  VIEW_SCREEN_MODES,
  FULLSCREEN_MODES_TO_DEFAULT_SCREEN_MODES,
} from '@/constants';

import { useActiveView } from '@/hooks/store/modules/active-view';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useFullscreen } from '@/hooks/fullscreen';
import { useView } from '@/hooks/store/modules/view';
import { useViewRouter } from '@/hooks/view/router';

import ViewShareLinkBtn from './view-share-link-btn.vue';
import ViewEditingBtn from './view-editing-btn.vue';
import ViewScrollTopBtn from './view-scroll-top-btn.vue';
import ViewScreenModeBtn from './view-screen-mode-btn.vue';
import ViewPeriodicRefreshBtn from './view-periodic-refresh-btn.vue';

export default {
  components: {
    ViewShareLinkBtn,
    ViewEditingBtn,
    ViewScrollTopBtn,
    ViewScreenModeBtn,
    ViewPeriodicRefreshBtn,
  },
  props: {
    activeTab: {
      type: Object,
      required: false,
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

      const element = document.querySelector('[data-app]');
      const viewElement = document.getElementById(`view-tab-${props.activeTab._id}`);

      if (!element) {
        return;
      }

      fullscreen.request(element, {
        fullscreenClass: 'full-screen',
        callback: (value) => {
          if (value) {
            viewElement.classList.add('view-fullscreen');
          } else {
            viewElement.classList.remove('view-fullscreen');
          }

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
      if (event.key === 'e' && event.ctrlKey && props.updatable) {
        toggleEditing();
        event.preventDefault();
        return;
      }

      /**
       * Check if Alt (Windows/Linux) or Command (Mac) + Shift + number keys are pressed
       */
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
            validationRules: 'required',
          },
          action: async (title) => {
            const data = {
              view: view.value._id,
              title,
            };

            await createViewTab({ data });
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

    .theme--dark & {
      border-color: #2196F3 !important;
      background-color: #2196F3 !important;
    }
  }

  &__add-edit-btn, &__add-fullscreen-btn  {
    border-color: #3f51b5 !important;
    background-color: #3f51b5 !important;

    .theme--dark & {
      border-color: #979797 !important;
      background-color: #979797 !important;
    }
  }

  &__fullscreen {
    z-index: 8;
  }
}

.view-fullscreen {
  overflow: auto;
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;

  z-index: 7;

  background: white;

  .theme--dark & {
    background: #424242;
  }
}
</style>
