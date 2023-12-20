<template lang="pug">
  v-card
    v-card-title.primary.pa-2.white--text
      v-layout(justify-space-between, align-center)
        h4 {{ $t('alarm.infoPopup') }}
        v-btn.ma-0.ml-3(
          color="white",
          icon,
          small,
          @click="$emit('close')"
        )
          v-icon(color="error", small) close
    v-fade-transition
      v-card-text.pa-2
        c-compiled-template(:template="template", :context="templateContext", :template-props="{ alarm }")
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
