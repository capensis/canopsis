<template lang="pug">
  v-navigation-drawer.side-bar.secondary(
  v-model="isOpen",
  :width="$config.SIDE_BAR_WIDTH",
  :class="{ editing: isEditingMode }"
  disable-resize-watcher,
  app,
  )
    div.brand.ma-0.secondary.lighten-1
      v-layout(justify-center, align-center)
        v-flex.text-xs-center(xs11)
          img.my-1(src="@/assets/canopsis.png")
        v-flex.version.white--text.caption
          div {{ version }}
    v-expansion-panel.panel(
    v-if="hasReadAnyViewAccess",
    expand,
    focusable,
    dark
    )
      v-expansion-panel-content.secondary.white--text(v-for="group in groups", :key="group._id")
        div.panel-header(slot="header")
          span(:title="group.name") {{ group.name }}
          v-btn(
          v-show="isEditingMode",
          depressed,
          small,
          icon,
          @click.stop="showEditGroupModal(group)"
          )
            v-icon(small) edit
        v-card.secondary.lighten-1.white--text(v-for="view in getAvailableViewsForGroup(group)", :key="view._id")
          router-link.panel-item-content-link(:title="view.title", :to="{ name: 'view', params: { id: view._id } }")
            v-card-text.panel-item-content
              v-layout(align-center, justify-space-between)
                v-flex
                  span.pl-2 {{ view.title }}
                v-flex
                  v-layout(justify-end)
                    v-btn.ma-0(
                    :v-show="checkViewEditButtonAccessById",
                    depressed,
                    small,
                    icon,
                    @click.prevent="showEditViewModal(view)"
                    )
                      v-icon(small) edit
                    v-btn.ma-0(
                    v-show="isEditingMode",
                    depressed,
                    small,
                    icon,
                    @click.prevent="showDuplicateViewModal(view)"
                    )
                      v-icon(small) file_copy
          v-divider
    v-divider
    groups-settings-button(
    :isEditingMode="isEditingMode",
    @toggleEditingMode="toggleEditingMode"
    )
</template>

<script>
import versionMixin from '@/mixins/entities/version';
import layoutNavigationGroupMenuMixin from '@/mixins/layout/navigation/group-menu';

import GroupsSettingsButton from './groups-settings-button.vue';

/**
 * Component for the side-bar, on the left of the application
 *
 * @prop {bool} [value=false] - visibility control
 *
 * @event input#update
 */
export default {
  components: { GroupsSettingsButton },
  mixins: [versionMixin, layoutNavigationGroupMenuMixin],
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isOpen: {
      get() {
        return this.value;
      },
      set(value) {
        if (value !== this.value) {
          this.$emit('input', value);
        }
      },
    },
    checkViewEditButtonAccessById() {
      return id => (this.checkUpdateViewAccessById(id) ||
        this.checkDeleteViewAccessById(this.view._id)) &&
        this.isEditingMode;
    },
  },
  mounted() {
    this.fetchVersion();
  },
};
</script>

<style lang="scss" scoped>
  a {
    color: inherit;
    text-decoration: none;
  }

  .panel {
    box-shadow: none;
  }

  .side-bar {
    position: fixed;
    height: 100vh;
    overflow-y: auto;

    & /deep/ .v-expansion-panel__header {
      height: 48px;
    }
  }

  .brand {
    height: 48px;
    position: relative;
  }

  .version {
    position: absolute;
    bottom: 0;
    right: 0;
    padding-right: 0.5em;
  }

  .panel-header {
    max-width: 88%;

    span {
      max-width: 100%;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      display: inline-block;
      vertical-align: middle;

      .editing & {
        max-width: 73%;
      }
    }
  }

  .panel-item-content {
    display: -webkit-box;
    display: -ms-flexbox;
    display: flex;
    cursor: pointer;
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    position: relative;
    padding: 12px 24px;
    height: 48px;

    & > div {
      max-width: 100%;
    }

    & /deep/ .v-btn:not(:last-child) {
      margin-right: 0;
    }

    .panel-item-content-link {
      max-width: 100%;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      display: inline-block;
      vertical-align: middle;
    }
  }
</style>
