<template lang="pug">
  div.secondary.darken-1.footer
    a(:href="$constants.CANOPSIS_DOCUMENTATION", target="_blank") {{ $t('login.documentation') }}
    a(:href="$constants.CANOPSIS_WEBSITE", target="_blank") {{ $t('login.website') }}
    a(:href="$constants.CANOPSIS_FORUM", target="_blank") {{ $t('login.forum') }}
    a.version(:href="changeLogHref", target="_blank") {{ version }}
</template>

<script>
import { CANOPSIS_DOCUMENTATION } from '@/constants';

import { entitiesInfoMixin } from '@/mixins/entities/info';

export default {
  mixins: [entitiesInfoMixin],
  computed: {
    changeLogHref() {
      if (this.version) {
        return `${CANOPSIS_DOCUMENTATION}/notes-de-version/${this.version}/`;
      }

      return CANOPSIS_DOCUMENTATION;
    },
  },
};
</script>

<style lang="scss">
  .footer {
    grid-area: footer;
    position: relative;
    color: white;
    min-height: 5em;
    margin-top: auto;
    display: flex;
    justify-content: center;
    align-items: center;

    a {
      color: inherit;
      text-decoration: none;
      padding: 0 2em;
    }

    .version {
      line-height: 5em;
      position: absolute;
      right: 0.5em;
      bottom: 0;
      text-decoration: underline;
      font-weight: bold;

      color: inherit;
      text-decoration: none;
    }
  }
</style>
