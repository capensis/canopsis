<template>
  <v-list class="comments-list">
    <v-list-item v-if="!comments || !comments.length">
      <v-list-item-title>{{ t('common.noData') }}</v-list-item-title>
    </v-list-item>

    <template v-for="(comment, index) in comments" :key="comment._id">
      <v-list-item>
        <v-list-item-title>{{ comment.author?.display_name }}</v-list-item-title>
        <v-list-item-subtitle>
          <c-compiled-template :template="comment.message" />
        </v-list-item-subtitle>
      </v-list-item>

      <v-divider v-if="index < comments.length - 1" />
    </template>
  </v-list>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    comments: {
      type: Array,
      default: () => [],
    },
  },

  setup() {
    const { t } = useI18n();

    return { t };
  },
};
</script>

<style lang="scss" scoped>
.comments-list {
  :deep(.v-list-item) {
    height: auto;
  }

  :deep(.v-list-item-subtitle) {
    word-break: break-word;
    white-space: pre-line;
  }
}
</style>
