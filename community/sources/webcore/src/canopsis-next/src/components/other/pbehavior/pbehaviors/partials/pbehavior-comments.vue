<template lang="pug">
  v-list.comments-list
    v-list-tile(v-if="!comments || !comments.length")
      v-list-tile-content
        v-list-tile-title {{ $t('common.noData') }}
    template(v-for="(comment, index) in comments")
      v-list-tile(:key="comment._id")
        v-list-tile-content
          v-list-tile-title {{ comment.author.name }}
          v-list-tile-sub-title {{ comment.message }}
      v-divider(v-if="index < comments.length - 1", :key="`divider-${index}`")
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
    & ::v-deep .v-list__tile {
      height: auto;

      &__sub-title {
        word-break: break-word;
        white-space: pre-line;
      }
    }
  }
</style>
