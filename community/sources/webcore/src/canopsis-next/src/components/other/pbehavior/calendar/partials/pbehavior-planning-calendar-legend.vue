<template>
  <v-menu
    :max-width="500"
    :nudge-width="500"
    :close-on-content-click="false"
    origin="left"
    left
    offset-x
    offset-overflow
  >
    <template #activator="{ on }">
      <v-tooltip top>
        <template #activator="{ on: tooltipOn }">
          <v-btn
            v-on="{ ...tooltipOn, ...on }"
            icon
          >
            <v-icon>info</v-icon>
          </v-btn>
        </template>
        <div>{{ $t('calendar.pbehaviorPlanningLegend.title') }}</div>
      </v-tooltip>
    </template>
    <v-card>
      <v-card-text>
        <template v-if="exceptionTypes.length">
          <div
            class="my-1"
            v-for="type in exceptionTypes"
            :key="type._id"
          >
            <div
              class="text-body-1"
              :style="getStyleForType(type)"
            >
              <v-icon
                class="px-1"
                color="white"
                small
              >
                {{ type.icon_name }}
              </v-icon>
              <strong>{{ type.name }}</strong>
            </div>
          </div>
        </template>
        <span
          class="text-subtitle-1"
          v-else
        >
          {{ $t('calendar.pbehaviorPlanningLegend.noData') }}
        </span>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script>
import { getMostReadableTextColor } from '@/helpers/color';

export default {
  props: {
    exceptionTypes: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    getStyleForType(type = {}) {
      return {
        backgroundColor: type.color,
        color: getMostReadableTextColor(type.color, { level: 'AA', size: 'large' }),
      };
    },
  },
};
</script>
