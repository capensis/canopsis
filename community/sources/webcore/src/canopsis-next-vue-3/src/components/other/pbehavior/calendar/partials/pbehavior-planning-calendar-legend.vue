<template>
  <v-menu :close-on-content-click="false" location="start" offset="8">
    <template #activator="{ props: menuProps }">
      <v-tooltip location="top">
        <template #activator="{ props: tooltipProps }">
          <v-btn v-bind="{ ...tooltipProps, ...menuProps }" variant="text" icon>
            <v-icon>info</v-icon>
          </v-btn>
        </template>

        <div>{{ $t('calendar.pbehaviorPlanningLegend.title') }}</div>
      </v-tooltip>
    </template>

    <v-card max-width="500">
      <v-card-text>
        <template v-if="exceptionTypes.length">
          <div v-for="type in exceptionTypes" :key="type._id" class="my-1">
            <div :style="getStyleForType(type)" class="text-body-1">
              <v-icon class="px-1" color="white" size="small">
                {{ type.icon_name }}
              </v-icon>

              <strong>{{ type.name }}</strong>
            </div>
          </div>
        </template>

        <span v-else class="text-subtitle-1">
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
