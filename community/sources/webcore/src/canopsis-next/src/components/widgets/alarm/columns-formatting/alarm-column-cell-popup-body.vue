<template>
  <v-card>
    <v-card-title class="primary pa-2 white--text">
      <v-layout
        class="gap-3"
        justify-space-between
        align-center
      >
        <h4>{{ $t('alarm.infoPopup') }}</h4>
        <v-btn
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
    <v-card-text class="pa-2">
      <c-compiled-template
        :template-id="templateId"
        :template="template"
        :context="templateContext"
        :template-props="{ alarm, selectedTag }"
        @select:tag="$emit('select:tag', $event)"
        @clear:tag="$emit('clear:tag')"
      />
    </v-card-text>
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
    templateId: {
      type: String,
      default: '',
    },
    selectedTag: {
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
