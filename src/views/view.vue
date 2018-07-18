<template lang="pug">
  div
    div(v-for="widgetWrapper in widgetWrappers", :key="widgetWrapper._id")
      div(:is="widgetsMap[widgetWrapper.widget.xtype]", :widget="widgetWrapper.widget")
</template>

<script>
import AlarmListContainer from '@/containers/alarm-list.vue';
import EntitiesListContainer from '@/containers/entities-list.vue';
import viewMixin from '@/mixins/view';

export default {
  components: {
    AlarmListContainer,
    EntitiesListContainer,
  },
  mixins: [
    viewMixin,
  ],
  props: {
    id: {
      type: [String, Number],
      required: true,
    },
  },
  data() {
    return {
      widgetsMap: {
        listalarm: 'alarm-list-container',
        crudcontext: 'entities-list-container',
      },
    };
  },
  mounted() {
    this.fetchView({ id: this.id });
  },
};
</script>
