<template>
  <c-alarm-actions-chips
    :class="{ 'my-1': !small }"
    :items="preparedLinks"
    :small="small"
    :inline-count="inlineCount"
    :item-class="chipsItemClass"
    class="c-links-chips"
    item-text="text"
    item-value="url"
    text-color=""
    return-object
    outlined
    @select="select"
    @activate="$emit('activate')"
  >
    <template #item="{ item }">
      <v-tooltip v-if="onlyIcon" location="top">
        <template #activator="{ props: activatorProps }">
          <v-icon v-bind="activatorProps" size="small">
            {{ item.icon }}
          </v-icon>
        </template>

        <span>{{ item.text }}</span>
      </v-tooltip>
      <c-layout v-else align="center">
        <v-icon class="mr-1" size="small">
          {{ item.icon }}
        </v-icon>
        <span>{{ item.text }}</span>
      </c-layout>
    </template>
  </c-alarm-actions-chips>
</template>

<script>
import { DEFAULT_LINKS_INLINE_COUNT, LINK_RULE_ACTIONS } from '@/constants';
import { writeTextToClipboard } from '@/helpers/clipboard';
import { harmonizeCategoryLinks, harmonizeLinks } from '@/helpers/entities/link/list';

export default {
  inject: ['$system'],

  emits: ['activate'],

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
        action: link.action,
        color: `blue-grey${this.$system.dark ? ' lighten-1' : ''}`,
        icon: link.icon_name,
        text: link.label,
        url: link.url,
      }));
    },

    chipsItemClass() {
      return {
        'c-links-chips__chip': true,
        'c-links-chips__chip--only-icon': this.onlyIcon,
      };
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
  },
};
</script>

<style lang="scss">
.c-links-chips__chip {
  .v-chip__content {
    padding: 0 4px;
  }

  &--only-icon {
    width: 26px;
  }
}
</style>
