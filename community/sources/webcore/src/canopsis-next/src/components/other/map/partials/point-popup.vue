<template lang="pug">
  v-card.point-popup(width="400")
    v-card-title.pa-2.white--text(:style="{ backgroundColor: color }")
      v-layout(justify-space-between, align-center)
        h4 {{ title }}
        v-btn.ma-0.ml-3(icon, small, @click="close")
          v-icon(color="white") close
    v-card-text
      c-compiled-template(
        v-if="point.entity && template",
        :template="template",
        :context="templateContext"
      )
      v-layout(v-else, column)
        span(v-if="point.entity") {{ $tc('common.entity') }}: {{ point.entity.name }}
        span(v-if="point.map") {{ $tc('common.map') }}: {{ point.map.name }}
    v-layout.ma-0.point-popup__actions(v-if="actions")
      v-btn.ma-0(
        v-if="hasAlarmsListAccess && point.entity",
        flat,
        block,
        @click.stop="$emit('show:alarms')"
      ) {{ $t('common.seeAlarms') }}
      v-btn.ma-0(
        v-if="point.map",
        flat,
        block,
        @click.stop="$emit('show:map')"
      )
        v-icon(left) link
        span.text-none  {{ point.map.name }}
</template>

<script>
import { isNumber } from 'lodash';

import { CSS_COLOR_VARS } from '@/config';
import { USERS_PERMISSIONS } from '@/constants';

import { getEntityColor } from '@/helpers/entities/entity/color';

import { authMixin } from '@/mixins/auth';

import MermaidPointMarker from './mermaid-point-marker.vue';

export default {
  components: { MermaidPointMarker },
  mixins: [authMixin],
  props: {
    point: {
      type: Object,
      required: true,
    },
    template: {
      type: String,
      required: false,
    },
    colorIndicator: {
      type: String,
      required: false,
    },
    actions: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    templateContext() {
      return { entity: this.point.entity };
    },

    color() {
      return isNumber(this.point.entity?.state)
        ? getEntityColor(this.point.entity, this.colorIndicator)
        : CSS_COLOR_VARS.primary;
    },

    title() {
      return this.point.entity ? this.point.entity.name : '';
    },

    hasAlarmsListAccess() {
      return this.checkAccess(USERS_PERMISSIONS.business.map.actions.alarmsList);
    },
  },
  methods: {
    close() {
      this.$emit('close');
    },
  },
};
</script>

<style lang="scss">
.point-popup {
  &__actions {
    background: #eee;

    .theme--dark & {
      background: #2f2f2f;
    }
  }
}
</style>
