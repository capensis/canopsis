<template lang="pug">
  v-tabs.view-tabs(
    ref="tabs",
    :key="vTabsKey",
    :value="$route.fullPath",
    :class="{ hidden: tabs.length < 2 && !editing, 'tabs-editing': editing }",
    :hide-slider="changed",
    color="secondary lighten-2",
    slider-color="primary",
    dark
  )
    c-draggable-list-field.d-flex(
      v-if="tabs.length",
      :value="tabs",
      :disabled="!editing",
      drag-class="draggable-item--dragging",
      chosen-class="draggable-item--chosen",
      @end="onDragEnd",
      @input="$emit('update:tabs', $event)"
    )
      v-tab.draggable-item(
        v-for="{ to, tab, title, key } in preparedTabs",
        :key="key",
        :disabled="changed",
        :to="to",
        exact,
        ripple
      )
        span {{ title }}
        template(v-if="updatable && editing")
          v-btn(small, flat, icon, @click.prevent="showUpdateTabModal(tab)")
            v-icon(small) edit
          v-btn(small, flat, icon, @click.prevent="showSelectViewModal(tab)")
            v-icon(small) file_copy
          v-btn(small, flat, icon, @click.prevent="showDeleteTabModal(tab)")
            v-icon(small) delete
    template(v-if="$scopedSlots.default")
      v-tabs-items(touchless)
        v-tab-item(
          v-for="{ to, tab, key } in preparedTabs",
          :key="key",
          :value="to",
          lazy
        )
          slot(:tab="tab", :visible="to === $route.fullPath")
</template>

<script>
import { MODALS, ROUTES_NAMES } from '@/constants';

import { activeViewMixin } from '@/mixins/active-view';
import { viewRouterMixin } from '@/mixins/view/router';
import { vuetifyTabsMixin } from '@/mixins/vuetify/tabs';
import { entitiesViewMixin } from '@/mixins/entities/view';
import { entitiesViewTabMixin } from '@/mixins/entities/view/tab';

export default {
  mixins: [
    activeViewMixin,
    viewRouterMixin,
    vuetifyTabsMixin,
    entitiesViewMixin,
    entitiesViewTabMixin,
  ],
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
  computed: {
    vTabsKey() {
      return this.view.tabs.map(tab => tab._id).join('-');
    },

    preparedTabs() {
      return this.tabs.map(tab => ({
        key: tab._id,
        to: this.getTabHrefById(tab._id),
        title: tab.title,
        tab,
      }));
    },

    getTabHrefById() {
      return (id) => {
        const { href } = this.$router.resolve({ query: { tabId: id } }, this.$route);

        return href.replace('#', '');
      };
    },
  },
  watch: {
    editing() {
      this.$nextTick(this.callTabsOnResizeMethod);
    },

    tabs: {
      immediate: true,
      handler() {
        this.onUpdateTabs();
      },
    },
  },
  methods: {
    showUpdateTabModal(tab) {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.viewTab.edit.title'),
          field: {
            name: 'text',
            label: this.$t('modals.viewTab.fields.title'),
            value: tab.title,
            validationRules: 'required',
          },
          action: title => this.updateViewTabAndFetch({
            id: tab._id,
            data: { ...tab, view: this.view._id, title },
          }),
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
            validationRules: 'required',
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

            this.$router.push({
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
        this.callTabsOnResizeMethod();
        this.callTabsUpdateTabsMethod();
      });
    },

    onDragEnd() {
      this.onUpdateTabs();
    },
  },
};
</script>

<style lang="scss" scoped>
  .view-tabs.hidden {
    & ::v-deep > .v-tabs__bar {
      display: none;
    }
  }

  .draggable-item {
    position: relative;
    transform: translateZ(0);

    .tabs-editing & {
      cursor: move;

      & ::v-deep .v-tabs__item {
        cursor: move;
      }
    }

    & ::v-deep .v-tabs__item--disabled {
      color: #fff;
      opacity: 1;

      button {
        color: rgba(255, 255, 255, 0.3) !important;
        box-shadow: none !important;
        pointer-events: none;
      }
    }

    &--dragging, &--chosen {
      & ::v-deep .v-ripple__container {
        display: none;
      }
    }

    &--dragging {
      background: var(--v-secondary-lighten2);
    }
  }
</style>
