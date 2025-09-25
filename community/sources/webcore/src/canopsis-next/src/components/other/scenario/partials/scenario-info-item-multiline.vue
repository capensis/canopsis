<template>
  <v-list-item class="text-subtitle-1">
    <v-list-item-content>
      <v-layout class="gap-2" column>
        <div class="text-subtitle-1">
          <span class="font-weight-bold">{{ label }}</span>
          <span v-if="value">:</span>
        </div>
        <v-layout v-if="value" class="gap-2 grey--text" align-center>
          <v-icon v-if="!hideIcon" color="grey">
            {{ icon }}
          </v-icon>
          <span v-if="value">{{ value }}</span>
        </v-layout>
      </v-layout>
    </v-list-item-content>
  </v-list-item>
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    value: {
      type: [String, Number],
      default: '',
    },
    hideIcon: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
