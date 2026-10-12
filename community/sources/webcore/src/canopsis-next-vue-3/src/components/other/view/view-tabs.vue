<template>
  <div class="view-tabs-wrapper">
    <v-tabs
      ref="tabs"
      :key="vTabsKey"
      :model-value="$route.query.tabId"
      :class="{
        hidden: tabs.length < 2 && !activeViewEditing,
        'tabs-editing': activeViewEditing,
      }"
      :hide-slider="changed"
      bg-color="secondary-lighten-2"
      class="view-tabs"
      color="white"
      slider-color="primary"
    >
      <c-draggable-list-field
        v-if="tabs.length"
        :model-value="tabs"
        :disabled="!activeViewEditing"
        class="d-flex"
        drag-class="draggable-item--dragging"
        chosen-class="draggable-item--chosen"
        @end="onDragEnd"
        @update:model-value="$emit('update:tabs', $event)"
      >
        <template #item="{ element: tab }">
          <v-tab
            :key="tab._id"
            :value="tab._id"
            :to="getTabRouteById(tab._id)"
            class="draggable-item view-tabs__tab"
            exact
            ripple
          >
            <span>{{ tab.title }}</span>

            <template v-if="updatable && activeViewEditing">
              <v-btn
                :disabled="changed"
                class="ml-2"
                size="small"
                variant="text"
                icon
                @click.prevent.stop="showUpdateTabModal(tab)"
              >
                <v-icon size="small">edit</v-icon>
              </v-btn>

              <v-btn
                :disabled="changed"
                size="small"
                variant="text"
                icon
                @click.prevent.stop="showSelectViewModal(tab)"
              >
                <v-icon size="small">file_copy</v-icon>
              </v-btn>

              <v-btn
                :disabled="changed"
                size="small"
                variant="text"
                icon
                @click.prevent.stop="showDeleteTabModal(tab)"
              >
                <v-icon size="small">delete</v-icon>
              </v-btn>
            </template>
          </v-tab>
        </template>
      </c-draggable-list-field>
    </v-tabs>

    <v-window v-if="$slots.default" :model-value="$route.query.tabId" touchless>
      <v-window-item v-for="tab in tabs" :key="tab._id" :value="tab._id">
        <slot :tab="tab" :visible="tab._id === $route.query.tabId" />
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { mapActions } from 'pinia';

import { MODALS, ROUTES_NAMES } from '@/constants';

import { useViewStore } from '@/stores/modules/view';

import { activeViewMixin } from '@/mixins/active-view';
import { viewRouterMixin } from '@/mixins/view/router';

export default {
  mixins: [activeViewMixin, viewRouterMixin],
  props: {
    tabs: {
      type: Array,
      required: true,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    changed: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:tabs'],
  computed: {
    vTabsKey() {
      return this.view.tabs.map(tab => tab._id).join('-');
    },

    getTabRouteById() {
      return id => ({
        query: {
          ...this.$route.query,
          tabId: id,
        },
      });
    },
  },
  watch: {
    activeViewEditing() {
      this.$nextTick(this.onUpdateTabs);
    },

    tabs: {
      immediate: true,
      handler() {
        this.onUpdateTabs();
      },
    },
  },
  methods: {
    ...mapActions(useViewStore, {
      updateViewTab: 'updateViewTab',
      copyViewTab: 'copyViewTab',
      removeViewTab: 'removeViewTab',
    }),

    showUpdateTabModal(tab) {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.viewTab.edit.title'),
          field: {
            name: 'text',
            label: this.$t('modals.viewTab.fields.title'),
            value: tab.title,
            validationRules: { required: true },
          },
          action: async (title) => {
            await this.updateViewTab({
              id: tab._id,
              data: { ...tab, view: this.view._id, title },
            });

            await this.fetchActiveView();
          },
        },
      });
    },

    showSelectViewModal(tab) {
      this.$modals.show({
        name: MODALS.selectView,
        config: {
          action: viewId => this.showCloneTabModal(tab, viewId),
        },
      });
    },

    showCloneTabModal(tab, viewId) {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.viewTab.duplicate.title'),
          field: {
            name: 'text',
            label: this.$t('modals.viewTab.fields.title'),
            validationRules: { required: true },
          },
          action: async (title) => {
            const data = {
              title,
              view: viewId,
            };

            const newTab = await this.copyViewTab({ id: tab._id, data });

            if (this.view._id === viewId) {
              await this.fetchActiveView();
            }

            await this.$router.push({
              name: ROUTES_NAMES.view,
              params: {
                id: viewId,
              },
              query: {
                tabId: newTab._id,
              },
            });
          },
        },
      });
    },

    showDeleteTabModal(tab) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeViewTab({ id: tab._id });
            await this.fetchActiveView();

            if (tab._id !== this.$route.query.tabId) {
              return;
            }

            if (this.view.tabs.length) {
              await this.redirectToFirstTab();
              return;
            }

            await this.redirectToViewRoot();
          },
        },
      });
    },

    onUpdateTabs() {
      this.$nextTick(() => {
        this.$refs.tabs?.updateSlider?.();
      });
    },

    onDragEnd() {
      this.onUpdateTabs();
    },
  },
};
</script>

<style lang="scss" scoped>
.view-tabs {
  &.hidden {
    display: none;
  }

  :deep(.v-tab:not(.v-tab--selected)) {
    color: rgba(255, 255, 255, 0.6);
  }

  :deep(.v-tab--selected) {
    color: #fff !important;
  }
}

.draggable-item {
  position: relative;
  transform: translateZ(0);

  .tabs-editing & {
    cursor: move;
  }

  &--dragging,
  &--chosen {
    :deep(.v-ripple__container) {
      display: none;
    }
  }

  &--dragging {
    background: var(--v-secondary-lighten2);
  }
}
</style>
