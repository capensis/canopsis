<template>
  <v-btn v-if="hasLlms" color="primary" @click="fixWithAi">
    <v-icon class="mr-2">
      $vuetify.icons.ai
    </v-icon>
    <span>{{ $t('pattern.fixJsonWithAi') }}</span>
  </v-btn>
</template>

<script>
import { computed, inject } from 'vue';

export default {
  props: {
    jsonString: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const aiChat = inject('$aiChat', {});

    const hasLlms = computed(() => aiChat.llms?.value?.length);

    const fixWithAi = () => aiChat.updateJsonString?.(props.jsonString);

    return {
      hasLlms,
      fixWithAi,
    };
  },
};
</script>
