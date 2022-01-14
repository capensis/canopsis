<template lang="pug">
  div.fab
    v-layout(row)
      view-scroll-top-btn
      view-periodic-refresh-btn
      v-speed-dial(
        v-if="updatable",
        v-model="isVSpeedDialOpen",
        direction="top",
        transition="slide-y-reverse-transition"
      )
        v-btn(
          slot="activator",
          :input-value="isVSpeedDialOpen",
          color="primary",
          dark,
          fab
        )
          v-icon menu
          v-icon close
        view-fullscreen-btn(:active-tab="activeTab", small, left-tooltip)
        view-editing-btn(v-if="updatable", :updatable="updatable")
        v-tooltip(left)
          v-btn(
            slot="activator",
            v-if="updatable",
            fab,
            dark,
            small,
            color="indigo",
            @click.stop="showCreateWidgetModal"
          )
            v-icon add
          span {{ $t('common.addWidget') }}
        v-tooltip(left)
          v-btn(
            slot="activator",
            v-if="updatable",
            fab,
            dark,
            small,
            color="green",
            @click.stop="showCreateTabModal"
          )
            v-icon add
          span {{ $t('common.addTab') }}
      view-fullscreen-btn(v-else, :active-tab="activeTab", top-tooltip)
</template>

<script>
import { MODALS } from '@/constants';

import { activeViewMixin } from '@/mixins/active-view';
import { entitiesViewTabMixin } from '@/mixins/entities/view/tab';

import ViewEditingBtn from '@/components/other/view/buttons/view-editing-btn.vue';
import ViewScrollTopBtn from '@/components/other/view/buttons/view-scroll-top-btn.vue';
import ViewFullscreenBtn from '@/components/other/view/buttons/view-fullscreen-btn.vue';
import ViewPeriodicRefreshBtn from '@/components/other/view/buttons/view-periodic-refresh-btn.vue';

export default {
  components: {
    ViewEditingBtn,
    ViewScrollTopBtn,
    ViewFullscreenBtn,
    ViewPeriodicRefreshBtn,
  },
  mixins: [
    activeViewMixin,
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
  },
  data() {
    return {
      isVSpeedDialOpen: false,
    };
  },

  methods: {
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

            return this.fetchActiveView();
          },
        },
      });
    },
  },
};
</script>
