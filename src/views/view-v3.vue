<template lang="pug">
  v-container
    div
      <!--div(v-for="widget in widgets", :key="widget._id")-->
        <!--h2 {{ widget.title }}-->
        <!--component(-->
        <!--:is="widgetsMap[widget.xtype]",-->
        <!--:widget="widget",-->
        <!--)-->
    v-speed-dial.fab(
    direction="left",
    :open-on-hover="true",
    transition="scale-transition"
    )
      v-btn(slot="activator", color="green darken-3", dark, fab)
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
  data() {
    return {
      widgetsMap: {
        /*        [WIDGET_TYPES.alarmList]: 'alarms-list',
        [WIDGET_TYPES.context]: 'entities-list',
        [WIDGET_TYPES.weather]: 'weather', */
      },
    };
  },
  mounted() {
    this.fetchView({ id: this.id });
  },
};
</script>
