<template>
  <v-list class="comments-list">
    <v-list-item v-if="!comments || !comments.length">
      <v-list-item-content>
        <v-list-item-title>{{ $t('common.noData') }}</v-list-item-title>
      </v-list-item-content>
    </v-list-item>
    <template v-for="(comment, index) in comments">
      <v-list-item :key="comment._id">
        <v-list-item-content>
          <v-list-item-title>{{ comment.author?.display_name }}</v-list-item-title>
          <v-list-item-subtitle>
            <c-compiled-template :template="comment.message" />
          </v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-divider
        v-if="index &lt; comments.length - 1"
        :key="`divider-${index}`"
      />
    </template>
  </v-list>
</template>

<script>
export default {
  props: {
    comments: {
      type: Array,
      default: () => ([]),
    },
  },
};
</script>

<style lang="scss" scoped>
  .comments-list {
    & ::v-deep .v-list-item {
      height: auto;

      &__sub-title {
        word-break: break-word;
        white-space: pre-line;
      }
    }
  }
</style>
