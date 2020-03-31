<template lang="pug">
  draggable(
    v-field="groups",
    :component-data="{ props: { expand: true, dark: true, focusable: true } }",
    :options="draggableOptions",
    element="v-expansion-panel"
  )
    group-panel(
      v-for="group in groups",
      :group="group",
      :key="group._id",
      hideActions
    )
      draggable-group-views(v-model="group.views")
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import DraggableGroupViews from './draggable-group-views.vue';
import GroupPanel from './group-panel.vue';

export default {
  components: { DraggableGroupViews, Draggable, GroupPanel },
  model: {
    prop: 'groups',
    event: 'change',
  },
  props: {
    groups: {
      type: Array,
      required: true,
    },
  },
  computed: {
    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        group: 'groups',
      };
    },
  },
};
</script>

