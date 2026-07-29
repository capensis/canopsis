<template>
  <v-card class="point-popup" width="400">
    <v-card-title :style="{ backgroundColor: color }" class="pa-2 text-white">
      <c-layout justify="space-between" align="center">
        <h4>{{ title }}</h4>
        <v-btn class="ml-3" color="white" icon="close" size="small" variant="text" @click="close" />
      </c-layout>
    </v-card-title>
    <v-card-text>
      <c-compiled-template
        v-if="point.entity && template"
        :template="template"
        :context="templateContext"
      />
      <c-layout v-else column>
        <span v-if="point.entity"> {{ $t('common.entity') }}: {{ point.entity.name }} </span>
        <span v-if="point.map"> {{ $t('common.map') }}: {{ point.map.name }} </span>
      </c-layout>
    </v-card-text>
    <c-layout v-if="actions" class="ma-0 bg-background">
      <v-btn
        v-if="hasAlarmsListAccess && point.entity"
        class="point-popup__action"
        variant="text"
        @click.stop="$emit('show:alarms')"
      >
        {{ $t('common.seeAlarms') }}
      </v-btn>
      <v-btn
        v-if="point.map"
        class="point-popup__action"
        variant="text"
        @click.stop="$emit('show:map')"
      >
        <v-icon start> link </v-icon>
        <span class="text-none text-truncate">
          {{ point.map.name }}
        </span>
      </v-btn>
    </c-layout>
  </v-card>
</template>

<script>
import { isNumber } from 'lodash';

import { CSS_COLORS_VARS } from '@/config';
import { USER_PERMISSIONS } from '@/constants';

import { getEntityColor } from '@/helpers/entities/entity/color';

import { authMixin } from '@/mixins/auth';

export default {
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

  emits: ['close', 'show:alarms', 'show:map'],

  computed: {
    templateContext() {
      return { entity: this.point.entity };
    },

    color() {
      return isNumber(this.point.entity?.state)
        ? getEntityColor(this.point.entity, this.colorIndicator)
        : CSS_COLORS_VARS.primary;
    },

    title() {
      return this.point.entity ? this.point.entity.name : '';
    },

    hasAlarmsListAccess() {
      return this.checkAccess(USER_PERMISSIONS.business.map.actions.alarmsList);
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
  &__action {
    width: 100%;
    flex-shrink: 1;

    .v-btn__content {
      max-width: 100%;
    }
  }
}
</style>
