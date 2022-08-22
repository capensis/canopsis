<template lang="pug">
  div
    mermaid-point-marker(
      v-for="point in points",
      :key="point._id",
      :x="point.x",
      :y="point.y",
      :entity="point.entity",
      :size="markerSize",
      :color-indicator="colorIndicator",
      :pbehavior-enabled="pbehaviorEnabled",
      ref="points",
      @click="openPopup(point, $event)"
    )
    v-menu(
      v-if="activePoint",
      :value="true",
      :position-x="positionX",
      :position-y="positionY",
      :close-on-content-click="false",
      ignore-click-outside,
      offset-overflow,
      offset-x,
      absolute,
      top
    )
      mermaid-point-popup(
        v-click-outside="clickOutsideDirective",
        :point="activePoint",
        :template="popupTemplate",
        :actions="popupActions",
        :color="getPointEntityColor(activePoint)",
        @show:alarms="showAlarmListModal",
        @show:map="showLinkedMap",
        @close="closePopup"
      )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import { generateDefaultAlarmListWidget } from '@/helpers/entities';
import { getEntityColor } from '@/helpers/color';

import { entitiesServiceEntityMixin } from '@/mixins/entities/service-entity';

import MermaidPointMarker from './mermaid-point-marker.vue';
import MermaidPointPopup from './mermaid-point-popup.vue';

const { mapActions: mapAlarmActions } = createNamespacedHelpers('alarm');

export default {
  components: { MermaidPointMarker, MermaidPointPopup },
  mixins: [entitiesServiceEntityMixin],
  props: {
    points: {
      type: Array,
      required: true,
    },
    markerSize: {
      type: Number,
      default: 24,
    },
    popupTemplate: {
      type: String,
      required: false,
    },
    popupActions: {
      type: Boolean,
      default: false,
    },
    alarmsColumns: {
      type: Array,
      required: false,
    },
    colorIndicator: {
      type: String,
      required: false,
    },
    pbehaviorEnabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      positionX: 0,
      positionY: 0,
      activePoint: undefined,
    };
  },
  computed: {
    clickOutsideDirective() {
      return {
        handler: this.closePopup,
        include: () => this.$refs.points.map(({ $el }) => $el),
        closeConditional: () => true,
      };
    },
  },
  methods: {
    ...mapAlarmActions({
      fetchComponentAlarmsListWithoutStore: 'fetchComponentAlarmsListWithoutStore',
    }),

    getPointEntityColor(point) {
      return point.entity ? getEntityColor(point.entity, this.colorIndicator) : undefined;
    },

    openPopup(point, event) {
      const { top, left, width } = event.target.getBoundingClientRect();

      this.positionY = top;
      this.positionX = left + width / 2;
      this.activePoint = point;
    },

    closePopup() {
      this.activePoint = undefined;
    },

    showLinkedMap() {
      this.$emit('show:map', this.activePoint.map);
      this.closePopup();
    },

    showAlarmListModal() {
      try {
        const widget = generateDefaultAlarmListWidget();

        widget.parameters.widgetColumns = this.alarmsColumns;

        this.$modals.show({
          name: MODALS.alarmsList,
          config: {
            widget,
            fetchList: params => this.fetchComponentAlarmsListWithoutStore({
              params: { ...params, _id: this.activePoint.entity._id },
            }),
          },
        });
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });
      }
    },
  },
};
</script>
