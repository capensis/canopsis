<template>
  <div class="view-fab-btns fab ma-2">
    <v-layout>
      <v-flex class="mr-3">
        <view-scroll-top-btn />
      </v-flex>
      <v-flex class="mr-3">
        <view-periodic-refresh-btn />
      </v-flex>
      <c-speed-dial
        direction="top"
        transition="slide-y-reverse-transition"
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
        <view-fullscreen-btn
          :value="fullscreen"
          :toggle-full-screen="toggleFullScreen"
          left-tooltip
          small
        />
        <view-editing-btn v-if="updatable" />
        <v-tooltip left>
          <template #activator="{ on }">
            <v-btn
              class="view-fab-btns__add-widget-btn"
              v-if="updatable"
              v-on="on"
              fab
              dark
              small
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
              v-on="on"
              color="green"
              fab
              dark
              small
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
import { MODALS } from '@/constants';

import { activeViewMixin } from '@/mixins/active-view';
import { viewRouterMixin } from '@/mixins/view/router';
import { entitiesViewTabMixin } from '@/mixins/entities/view/tab';

import ViewShareLinkBtn from './view-share-link-btn.vue';
import ViewEditingBtn from './view-editing-btn.vue';
import ViewScrollTopBtn from './view-scroll-top-btn.vue';
import ViewFullscreenBtn from './view-fullscreen-btn.vue';
import ViewPeriodicRefreshBtn from './view-periodic-refresh-btn.vue';

export default {
  components: {
    ViewShareLinkBtn,
    ViewEditingBtn,
    ViewScrollTopBtn,
    ViewFullscreenBtn,
    ViewPeriodicRefreshBtn,
  },
  mixins: [
    activeViewMixin,
    viewRouterMixin,
    entitiesViewTabMixin,
  ],
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
  data() {
    return {
      opened: false,
      fullscreen: false,
    };
  },
  created() {
    document.addEventListener('keydown', this.keyDownListener);
  },
  beforeDestroy() {
    this.$fullscreen.exit();
    document.removeEventListener('keydown', this.keyDownListener);
  },
  methods: {
    toggleFullScreen() {
      if (!this.activeTab) {
        this.$popups.warning({ text: this.$t('view.errors.emptyTabs') });
        return;
      }

      const element = document.querySelector('[data-app]');
      const viewElement = document.getElementById(`view-tab-${this.activeTab._id}`);

      if (!element) {
        return;
      }

      this.$fullscreen.toggle(element, {
        fullscreenClass: 'full-screen',
        callback: (value) => {
          if (value) {
            viewElement.classList.add('view-fullscreen');
          } else {
            viewElement.classList.remove('view-fullscreen');
          }

          this.fullscreen = value;
        },
      });
    },

    keyDownListener(event) {
      if (event.key === 'e' && event.ctrlKey && this.updatable) {
        this.toggleEditing();
        event.preventDefault();
      } else if (event.key === 'Enter' && event.altKey) {
        this.toggleFullScreen();
        event.preventDefault();
      }
    },

    showCreateWidgetModal() {
      if (!this.activeTab) {
        this.$popups.warning({ text: this.$t('view.errors.emptyTabs') });
        return;
      }

      this.$modals.show({
        name: MODALS.createWidget,
        config: {
          tab: this.activeTab,
        },
      });
    },

    showCreateTabModal() {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.viewTab.create.title'),
          field: {
            name: 'text',
            label: this.$t('modals.viewTab.fields.title'),
            validationRules: 'required',
          },
          action: async (title) => {
            const data = {
              view: this.view._id,
              title,
            };

            await this.createViewTab({ data });
            await this.fetchActiveView();

            if (!this.$route.query.tabId) {
              await this.redirectToFirstTab();
            }
          },
        },
      });
    },
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
