<template lang="pug">
  c-alarm-actions-chips.c-alarm-links-chips(
    v-if="hasAccessToLinks",
    :items="links",
    :small="small",
    :inline-count="inlineCount",
    item-text="text",
    item-value="url",
    item-class="c-alarm-links-chips__chip",
    return-object,
    @select="select",
    @activate="activate"
  )
    template(#item="{ item }")
      v-tooltip(v-if="onlyIcon", top, custom-activator)
        template(#activator="{ on }")
          v-icon(v-on="on", color="white", small) {{ item.icon }}
        span {{ item.text }}
      template(v-else)
        v-icon(color="white", small) {{ item.icon }}
        span.ml-1 {{ item.text }}
</template>

<script>
import {
  ALARM_LIST_ACTIONS_TYPES,
  DEFAULT_LINKS_INLINE_COUNT,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  LINK_RULE_ACTIONS,
} from '@/constants';

import { harmonizeLinks, harmonizeCategoryLinks } from '@/helpers/links';

import { authMixin } from '@/mixins/auth';

export default {
  mixins: [authMixin],
  props: {
    alarm: {
      type: Object,
      required: true,
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
    hasAccessToLinks() {
      return this.checkAccess(BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[ALARM_LIST_ACTIONS_TYPES.links]);
    },

    links() {
      const links = this.category
        ? harmonizeCategoryLinks(this.alarm.links, this.category)
        : harmonizeLinks(this.alarm.links);

      return links.map(link => ({
        text: link.label,
        icon: link.icon_name,
        url: link.url,
        color: 'grey',
      }));
    },
  },
  methods: {
    async select(link) {
      if (link.action === LINK_RULE_ACTIONS.copy) {
        try {
          await navigator.clipboard.writeText(link.url);

          this.$popups.success({ text: this.$t('testSuite.popups.systemMessageCopied') });
        } catch (err) {
          console.error(err);

          this.$popups.error({ text: this.$t('errors.default') });
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
.c-alarm-links-chips__chip .v-chip__content {
  padding: 0 4px;
}
</style>
