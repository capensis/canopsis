<template>
  <v-layout
    class="ai-chat-messages-list gap-3"
    column
  >
    <template v-for="message in messages">
      <ai-chat-message
        :key="message._id"
        :message="message"
        :history="history"
      />
      <ai-chat-pattern
        v-if="message.patterns"
        :key="`${message._id}-pattern`"
        :version="message.version"
        :active="message.version === activeVersion"
        :patterns="message.patterns"
        :from-version="message.from_version"
        :history="history"
        @restore:version="restoreVersion"
      />
    </template>
    <ai-chat-message
      v-if="thinkingMessage"
      :key="thinkingMessage._id"
      :message="thinkingMessage"
    />
  </v-layout>
</template>

<script>
import AiChatMessage from './partials/ai-chat-message.vue';
import AiChatPattern from './partials/ai-chat-pattern.vue';

export default {
  components: {
    AiChatMessage,
    AiChatPattern,
  },
  props: {
    messages: {
      type: Array,
      required: true,
    },
    thinkingMessage: {
      type: Object,
      default: null,
    },
    activeVersion: {
      type: Number,
      default: null,
    },
    history: {
      type: Boolean,
      default: false,
    },
  },
  setup(_, { emit }) {
    const restoreVersion = version => emit('restore:version', version);

    return {
      restoreVersion,
    };
  },
};
</script>

<style lang="scss" scoped>
.ai-chat-messages-list {
  flex: 0 0 auto !important;
  min-height: 0;

  > * {
    flex: 0 0 auto !important;
  }
}
</style>
