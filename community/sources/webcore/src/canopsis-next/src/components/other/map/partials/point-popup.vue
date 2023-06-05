<template lang="pug">
  v-card.point-popup(width="400")
    v-card-title.pa-2.white--text(:style="{ backgroundColor: color }")
      v-layout(justify-space-between, align-center)
        h4 {{ title }}
        v-btn.ma-0.ml-3(icon, small, @click="close")
          v-icon(color="white") close
    v-card-text
      v-runtime-template(
        v-if="point.entity && template",
        :template="compiledTemplate"
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
      ) {{ $t('serviceWeather.seeAlarms') }}
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
import VRuntimeTemplate from 'v-runtime-template';

import { COLORS } from '@/config';

import { USERS_PERMISSIONS } from '@/constants';

import { compile } from '@/helpers/handlebars';
import { getEntityColor } from '@/helpers/color';

import { authMixin } from '@/mixins/auth';

import MermaidPointMarker from './mermaid-point-marker.vue';

export default {
  components: { VRuntimeTemplate, MermaidPointMarker },
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
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.template, { entity: this.point.entity });

        return `<div>${compiledTemplate}</div>`;
      },
      lazy: true,
      default: '',
    },
  },
  computed: {
    color() {
      return isNumber(this.point.entity?.state)
        ? getEntityColor(this.point.entity, this.colorIndicator)
        : COLORS.primary;
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
