<template lang="pug">
  v-container
    div
      div(v-for="widget in widgets", :key="widget._id")
        h2 {{ widget.title }}
        div(
        :is="widgetsMap[widget.xtype]",
        :widget="widget",
        @openSettings="openSettings(widget)"
        )
    v-speed-dial.fab(
    direction="top",
    :open-on-hover="true",
    transition="scale-transition"
    )
      v-btn(slot="activator", v-model="fab", color="green darken-3", dark, fab)
        v-icon add
      v-tooltip(left)
        v-btn(slot="activator", fab, dark, small, color="indigo", @click.prevent="showCreateWidgetModal")
          v-icon widgets
        span {{ $t('common.widget') }}
</template>

<script>
import entitiesViewV3Mixin from '@/mixins/entities/view-v3/view-v3';

export default {
  mixins: [entitiesViewV3Mixin],
  props: {
    id: {
      type: [String, Number],
      required: true,
    },
  },
  mounted() {
    this.fetchView({ id: this.id });
  },
};
</script>
