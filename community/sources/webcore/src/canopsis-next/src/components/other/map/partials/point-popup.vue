<template>
  <v-card
    class="point-popup"
    width="400"
  >
    <v-card-title
      class="pa-2 white--text"
      :style="{ backgroundColor: color }"
    >
      <v-layout
        justify-space-between
        align-center
      >
        <h4>{{ title }}</h4>
        <v-btn
          class="ma-0 ml-3"
          icon
          small
          @click="close"
        >
          <v-icon color="white">
            close
          </v-icon>
        </v-btn>
      </v-layout>
    </v-card-title>
    <v-card-text>
      <c-compiled-template
        v-if="point.entity && template"
        :template="template"
        :context="templateContext"
      />
      <v-layout
        v-else
        column
      >
        <span v-if="point.entity">
          {{ $tc('common.entity') }}: {{ point.entity.name }}
        </span>
        <span v-if="point.map">{{ $tc('common.map') }}: {{ point.map.name }}</span>
      </v-layout>
    </v-card-text>
    <v-layout
      class="ma-0 background darken-1"
      v-if="actions"
    >
      <v-btn
        class="point-popup__action"
        v-if="hasAlarmsListAccess && point.entity"
        text
        @click.stop="$emit('show:alarms')"
      >
        {{ $t('common.seeAlarms') }}
      </v-btn>
      <v-btn
        class="point-popup__action"
        v-if="point.map"
        text
        @click.stop="$emit('show:map')"
      >
        <v-icon left>
          link
        </v-icon>
        <span class="text-none text-truncate">{{ point.map.name }}</span>
      </v-btn>
    </v-layout>
  </v-card>
</template>

<script>
import { isNumber } from 'lodash';

import { CSS_COLORS_VARS } from '@/config';
import { USERS_PERMISSIONS } from '@/constants';

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
  &__action {
    width: 100%;
    flex-shrink: 1;

    .v-btn__content {
      max-width: 100%;
    }
  }
}
</style>
