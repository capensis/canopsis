<template lang="pug">
  div.view(:id="`view-tab-${tab._id}`")
    v-layout(v-for="row in rows", :key="row._id", wrap)
      v-flex(xs12)
        v-layout.hide-on-full-screen(justify-end)
          v-btn.ma-2(
            data-test="deleteRowButton",
            v-if="isEditingMode && hasUpdateAccess",
            @click.stop="showDeleteRowModal(row)",
            small,
            color="error"
          ) {{ $t('view.deleteRow') }} - {{ row.title }}
      v-flex(
        v-for="widget in row.widgets",
        :key="widget._id",
        :class="getWidgetFlexClass(widget)"
      )
        widget-wrapper(
          :widget="widget",
          :tab="tab",
          :isEditingMode="isEditingMode",
          :row="row",
          :updateTabMethod="updateTabMethod"
        )
</template>

<script>
import { MODALS } from '@/constants';

import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';

import popupMixin from '@/mixins/popup';
import modalMixin from '@/mixins/modal';
import sideBarMixin from '@/mixins/side-bar/side-bar';

export default {
  components: {
    WidgetWrapper,
  },
  mixins: [
    popupMixin,
    modalMixin,
    sideBarMixin,
  ],
  props: {
    tab: {
      type: Object,
      required: true,
    },
    hasUpdateAccess: {
      type: Boolean,
      default: false,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
    updateTabMethod: {
      type: Function,
      required: true,
    },
  },
  computed: {
    rows() {
      return this.tab.rows || [];
    },

    getWidgetFlexClass() {
      return widget => [
        `xs${widget.size.sm}`,
        `md${widget.size.md}`,
        `lg${widget.size.lg}`,
      ];
    },
  },
  methods: {
    showDeleteRowModal(row = {}) {
      const widgets = row.widgets || [];

      if (widgets.length > 0) {
        this.addErrorPopup({ text: this.$t('errors.lineNotEmpty') });
      } else {
        this.showModal({
          name: MODALS.confirmation,
          config: {
            action: () => {
              const newTab = {
                ...this.tab,

                rows: this.rows.filter(tabRow => tabRow._id !== row._id),
              };

              return this.updateTabMethod(newTab);
            },
          },
        });
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .full-screen {
    .hide-on-full-screen {
      display: none;
    }
  }
</style>
