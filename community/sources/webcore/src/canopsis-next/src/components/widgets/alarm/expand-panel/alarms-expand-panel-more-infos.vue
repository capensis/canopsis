<template>
  <div>
    <c-compiled-template
      v-if="template"
      :template-id="templateId"
      :template="template"
      :context="templateContext"
      @select:tag="$emit('select:tag', $event)"
      @remove:tag="$emit('remove:tag', $event)"
    />

    <c-layout
      v-else
      justify="center"
      align="center"
      gap="2"
    >
      <v-icon color="info">
        infos
      </v-icon>
      <p class="ma-0">
        {{ $t('alarm.moreInfos.defineATemplate') }}
      </p>
    </c-layout>
  </div>
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
    selectedTags: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['select:tag', 'remove:tag'],

  computed: {
    templateContext() {
      return {
        alarm: this.alarm,
        entity: this.alarm.entity,
      };
    },
  },
};
</script>
