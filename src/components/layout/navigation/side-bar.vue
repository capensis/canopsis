<template lang="pug">
  v-navigation-drawer.side-bar.grey.darken-4(
  v-model="isOpen",
  :clipped="$mq | mq({ m: false, l: true })",
  :width="width",
  absolute,
  app,
  )
    v-expansion-panel(
    class="panel",
    expand,
    focusable,
    dark
    )
      v-expansion-panel-content(v-for="group in groups", :key="group._id").grey.darken-4.white--text
        div(slot="header") {{ group.name }}
        v-card.grey.darken-3.white--text(v-for="view in group.views", :key="view._id")
          v-card-text
            router-link(:to="{ name: 'view', params: { id: view._id } }") {{ view.title }}
    v-divider
    v-btn.addBtn(
    fab,
    dark,
    fixed,
    bottom,
    right,
    color="green darken-4",
    @click="showCreateViewModal"
    )
      v-icon(dark) add
</template>

<script>
import { SIDE_BAR_WIDTH } from '@/config';
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
      width: SIDE_BAR_WIDTH,
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
    showCreateViewModal() {
      this.showModal({
        name: MODALS.createView,
      });
    },
  },
};
</script>

<style scoped>
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
