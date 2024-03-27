<template>
  <c-alarm-actions-chips
    :class="{ 'my-1': !small }"
    :items="preparedLinks"
    :small="small"
    :inline-count="inlineCount"
    class="c-links-chips"
    item-text="text"
    item-value="url"
    item-class="c-links-chips__chip"
    text-color=""
    return-object
    outlined
    @select="select"
    @activate="activate"
  >
    <template #item="{ item }">
      <v-tooltip
        v-if="onlyIcon"
        top
      >
        <template #activator="{ on }">
          <v-icon
            small
            v-on="on"
          >
            {{ item.icon }}
          </v-icon>
        </template>
        <span>{{ item.text }}</span>
      </v-tooltip>
      <v-layout
        v-else
        align-center
      >
        <v-icon
          class="mr-1"
          small
        >
          {{ item.icon }}
        </v-icon>
        <span>{{ item.text }}</span>
      </v-layout>
    </template>
  </c-alarm-actions-chips>
</template>

<script>
import { DEFAULT_LINKS_INLINE_COUNT, LINK_RULE_ACTIONS } from '@/constants';

import { harmonizeLinks, harmonizeCategoryLinks } from '@/helpers/entities/link/list';
import { writeTextToClipboard } from '@/helpers/clipboard';

export default {
  inject: ['$system'],
  props: {
    links: {
      type: Object,
      required: false,
    },
    inlineCount: {
      type: [Number, String],
      default: DEFAULT_LINKS_INLINE_COUNT,
    },
    small: {
      type: Boolean,
      default: false,
    },
    category: {
      type: String,
      required: false,
    },
    onlyIcon: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    preparedLinks() {
      const links = this.category
        ? harmonizeCategoryLinks(this.links, this.category)
        : harmonizeLinks(this.links);

      return links.map(link => ({
        text: link.label,
        icon: link.icon_name,
        url: link.url,
        action: link.action,
        color: `blue-grey${this.$system.dark ? ' lighten-1' : ''}`,
      }));
    },
  },
  methods: {
    async select(link) {
      if (link.action === LINK_RULE_ACTIONS.copy) {
        try {
          await writeTextToClipboard(link.url);

          this.$popups.success({ text: this.$t('popups.copySuccess') });
        } catch (err) {
          console.error(err);

          this.$popups.error({ text: this.$t('popups.copyError') });
        }

        return;
      }

      window.open(link.url, '_blank');
    },

    activate() {
      this.$emit('activate');
    },
  },
};
</script>

<style lang="scss">
.c-links-chips__chip .v-chip__content {
  padding: 0 4px;
}
</style>
