<template lang="pug">
  v-navigation-drawer.side-bar.secondary(
  v-model="isOpen",
  :width="width",
  app,
  )
    div.brand.ma-0.secondary.lighten-1
      v-layout(justify-center, align-center)
        img.my-1(src="@/assets/canopsis.png")
    v-expansion-panel.panel(
    expand,
    focusable,
    dark
    )
      v-expansion-panel-content(v-for="group in groups", :key="group._id").secondary.white--text
        div(slot="header")
          span {{ group.name }}
          v-btn(
          v-show="editing",
          depressed,
          small,
          icon,
          @click.stop="showEditGroupModal(group)"
          )
            v-icon(small) edit
        v-card.secondary.white--text(v-for="view in group.views", :key="view._id")
          v-card-text
            router-link(:to="{ name: 'view', params: { id: view._id } }")
              span.pl-3 {{ view.title }}
              v-btn(
              v-show="editing",
              depressed,
              small,
              icon,
              color="grey darken-2",
              @click.prevent="showEditViewModal(view)"
              )
                v-icon(small) edit
              v-btn(
              v-show="editing",
              depressed,
              small,
              icon,
              color="grey darken-2",
              @click.prevent="showDuplicateViewModal(view)"
              )
                v-icon(small) file_copy
    v-divider
    v-speed-dial(
    v-model="fab"
    bottom,
    right,
    direction="top"
    transition="slide-y-reverse-transition"
    )
      v-tooltip(slot="activator", right)
        v-btn.primary(slot="activator", v-model="fab", fab, dark)
          v-icon settings
          v-icon close
        span {{ $t('layout.sideBar.buttons.settings') }}
      v-tooltip(right)
        v-btn(
        slot="activator",
        v-model="editing",
        fab,
        dark,
        small,
        color="blue darken-4",
        @click.stop="editModeToggle"
        )
          v-icon(dark) edit
          v-icon(dark) done
        span {{ $t('layout.sideBar.buttons.edit') }}
      v-tooltip(right)
        v-btn(
        slot="activator",
        fab,
        dark,
        small,
        color="green darken-4",
        @click.stop="showCreateViewModal"
        )
          v-icon(dark) add
        span {{ $t('layout.sideBar.buttons.create') }}
</template>

<script>
import { MODALS } from '@/constants';
import modalMixin from '@/mixins/modal/modal';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';

/**
 * Component for the side-bar, on the left of the application
 *
 * @prop {bool} [value=false] - visibility control
 *
 * @event input#update
 */
export default {
  mixins: [
    modalMixin,
    entitiesViewGroupMixin,
  ],
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      fab: false,
      editing: false,
      width: this.$config.SIDE_BAR_WIDTH,
    };
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
  },
  mounted() {
    this.fetchGroupsList();
  },
  methods: {
    editModeToggle() {
      this.editing = !this.editing;
    },

    showEditGroupModal(group) {
      this.showModal({
        name: MODALS.createGroup,
        config: { group },
      });
    },

    showEditViewModal(view) {
      this.showModal({
        name: MODALS.createView,
        config: { view },
      });
    },

    showCreateViewModal() {
      this.showModal({
        name: MODALS.createView,
      });
    },
    showDuplicateViewModal(view) {
      this.showModal({
        name: MODALS.createView,
        config: {
          view,
          isDuplicating: true,
        },
      });
    },
  },
};
</script>

<style scoped>
  .v-speed-dial {
    position: fixed;
  }

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
  }
</style>
