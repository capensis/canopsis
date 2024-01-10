<template>
  <v-card>
    <v-card-title class="primary pa-2 white--text">
      <v-layout
        justify-space-between
        align-center
      >
        <h4>{{ $t('alarm.infoPopup') }}</h4>
        <v-btn
          class="ma-0 ml-3"
          color="white"
          icon
          small
          @click="$emit('close')"
        >
          <v-icon
            color="error"
            small
          >
            close
          </v-icon>
        </v-btn>
      </v-layout>
    </v-card-title>
    <v-fade-transition>
      <v-card-text class="pa-2">
        <c-compiled-template
          :template="template"
          :context="templateContext"
          :template-props="{ alarm }"
        />
      </v-card-text>
    </v-fade-transition>
  </v-card>
</template>

<script>
export default {
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    template: {
      type: String,
      default: '',
    },
  },
  computed: {
    templateContext() {
      return {
        alarm: this.alarm,
        entity: this.alarm.entity ?? {},
      };
    },
  },
};
</script>
