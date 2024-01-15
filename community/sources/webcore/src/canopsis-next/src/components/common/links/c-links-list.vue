<template>
  <div class="mt-1">
    <div
      v-for="(categoryLinks, categoryName) in preparedLinks"
      :key="categoryName"
    >
      <span class="category mr-2">{{ categoryName }}</span>
      <v-divider light />
      <div
        v-for="(link, index) in categoryLinks"
        :key="index"
      >
        <div class="pa-2 text-right">
          <c-copy-wrapper
            v-if="link.action === $constants.LINK_RULE_ACTIONS.copy"
            :value="link.url"
          >
            {{ link.label }}
          </c-copy-wrapper>
          <a
            v-else
            :href="link.url"
            target="_blank"
          >
            {{ link.label }}
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { harmonizeCategoryLinks, harmonizeCategoriesLinks } from '@/helpers/entities/link/list';

export default {
  props: {
    links: {
      type: Object,
      default: () => ({}),
    },
    category: {
      type: String,
      default: null,
    },
  },
  computed: {
    preparedLinks() {
      return this.category
        ? { [this.category]: harmonizeCategoryLinks(this.links, this.category) }
        : harmonizeCategoriesLinks(this.links);
    },
  },
};
</script>

<style lang="scss" scoped>
.category {
  display: inline-block;

  &:first-letter {
    text-transform: uppercase;
  }
}
</style>
