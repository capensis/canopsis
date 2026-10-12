<template>
  <v-card class="alarm-column-cell-popup-body">
    <v-card-title class="alarm-column-cell-popup-body__title bg-primary pa-2 text-white">
      <c-layout class="gap-3" justify="space-between" align="center">
        <span class="alarm-column-cell-popup-body__title-text">
          {{ $t('alarm.infoPopup') }}
        </span>
        <v-btn
          color="error"
          density="compact"
          size="x-small"
          variant="text"
          icon
          @click="$emit('close')"
        >
          <v-icon color="error" size="small">
            close
          </v-icon>
        </v-btn>
      </c-layout>
    </v-card-title>
    <v-card-text class="pa-2">
      <c-compiled-template
        :template-id="templateId"
        :template="template"
        :context="templateContext"
        :template-props="templateProps"
        @select:tag="$emit('select:tag', $event)"
        @remove:tag="$emit('remove:tag', $event)"
      />
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  emits: ['close', 'select:tag', 'remove:tag'],

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
    selectedTags: {
      type: Array,
      default: () => [],
    },
  },

  computed: {
    templateProps() {
      return {
        alarm: this.alarm,
        selectedTags: this.selectedTags,
      };
    },

    templateContext() {
      return {
        alarm: this.alarm,
        entity: this.alarm.entity ?? {},
      };
    },
  },
};
</script>

<style lang="scss">
.alarm-column-cell-popup-body {
  &__title {
    min-height: 32px;
  }

  &__title-text {
    font-size: 14px;
    font-weight: 500;
    line-height: 20px;
  }
}
</style>
